package engine

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/engine/actions"
	"github.com/fedulovivan/mhz19-go/internal/engine/conditions"
	"github.com/fedulovivan/mhz19-go/internal/message_queue"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var queuesContainer = message_queue.NewContainer()

type GetProviderFn func(ch types.ChannelType) types.ChannelProvider

var tidSeq = utils.NewSeq()

type engine struct {
	logTag         types.LogTagFn
	providers      []types.ChannelProvider
	rules          []types.Rule
	messageService types.MessagesService
	devicesService types.DevicesService
	ldmService     types.LdmService
}

func NewEngine() types.Engine {
	return &engine{
		logTag: func(m string) string { return m },
	}
}

func (e *engine) SetLogTag(f types.LogTagFn) {
	e.logTag = f
}
func (e *engine) SetProviders(s ...types.ChannelProvider) {
	e.providers = s
}
func (e *engine) SetMessagesService(s types.MessagesService) {
	e.messageService = s
}
func (e *engine) SetDevicesService(s types.DevicesService) {
	e.devicesService = s
}
func (e *engine) SetLdmService(r types.LdmService) {
	e.ldmService = r
}
func (e *engine) AppendRules(rules ...types.Rule) {
	slog.Debug(e.logTag("AppendRules"), "len", len(rules))
	e.rules = append(e.rules, rules...)
}

//	func (e *engine) LogTag(in string) string {
//		return e.logTag(in)
//	}

func (e *engine) MessagesService() types.MessagesService {
	return e.messageService
}
func (e *engine) DevicesService() types.DevicesService {
	return e.devicesService
}

func (e *engine) Start() {
	for _, p := range e.providers {
		go func(provider types.ChannelProvider) {
			provider.Init()
			for message := range provider.MessageChan() {
				e.HandleMessage(message, e.rules)
			}
		}(p)
	}
}

func (e *engine) Provider(ct types.ChannelType) types.ChannelProvider {
	for _, provider := range e.providers {
		if provider.Channel() == ct {
			return provider
		}
	}
	return nil
}

func (e *engine) Stop() {
	for _, s := range e.providers {
		s.Stop()
	}
}

func (e *engine) InvokeConditionFunc(mt types.MessageTuple, fn types.CondFn, args types.Args, r types.Rule, tid string) bool {
	impl := conditions.Get(fn)
	res := impl(mt, args, e)
	slog.Debug(e.logTag(tid+fmt.Sprintf("Rule #%v condition exec", r.Id)), "fn", fn, "args", args, "res", res)
	return res
}

func (e *engine) InvokeActionFunc(mm []types.Message, a types.Action, r types.Rule, tid string) {
	impl := actions.Get(a.Fn)
	slog.Debug(e.logTag(tid+fmt.Sprintf("Rule #%v action exec", r.Id)), "fn", a.Fn, "args", a.Args)
	go func() {
		err := impl(mm, a, e)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error: %s", a.Fn, err))
		}
	}()
}

func (e *engine) MatchesListSome(mt types.MessageTuple, cc []types.Condition, r types.Rule, tid string) bool {
	for _, c := range cc {
		if e.MatchesCondition(mt, c, r, tid) {
			return true
		}
	}
	return false
}

func (e *engine) MatchesListEvery(mt types.MessageTuple, cc []types.Condition, r types.Rule, tid string) bool {
	if len(cc) == 0 {
		return false
	}
	for _, c := range cc {
		if !e.MatchesCondition(mt, c, r, tid) {
			return false
		}
	}
	return true
}

func (e *engine) MatchesCondition(mt types.MessageTuple, c types.Condition, r types.Rule, tid string) bool {
	withFn := c.Fn != 0
	withList := len(c.List) > 0
	if withFn && !withList {
		return e.InvokeConditionFunc(mt, c.Fn, c.Args, r, tid)
	} else if withList && !withFn {
		if c.Or {
			return e.MatchesListSome(mt, c.List, r, tid)
		} else {
			return e.MatchesListEvery(mt, c.List, r, tid)
		}
	} else {
		return true
	}
}

func (e *engine) ExecuteActions(mm []types.Message, r types.Rule, tid string) {
	slog.Debug(e.logTag(tid + fmt.Sprintf("Rule #%v going to execute %v actions", r.Id, len(r.Actions))))
	for _, a := range r.Actions {
		e.InvokeActionFunc(mm, a, r, tid)
	}
}

// called simultaneusly upon receiving messages from all providers
// should be thread-safe
func (e *engine) HandleMessage(m types.Message, rules []types.Rule) {
	tid := fmt.Sprintf("Tid #%v ", tidSeq.Next())
	p := m.Payload
	if m.DeviceClass == types.DEVICE_CLASS_ZIGBEE_BRIDGE {
		p = "<too big to render>"
	}
	slog.Debug(
		e.logTag(tid+"New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)
	ldmKey := e.ldmService.MakeKey(m.DeviceClass, m.DeviceId)
	// var prev types.Message
	tuple := types.MessageTuple{
		Curr: &m,
	}
	if e.ldmService.Has(ldmKey) {
		prev := e.ldmService.Get(ldmKey)
		tuple.Prev = &prev
	}
	slog.Debug(e.logTag(tid + fmt.Sprintf("Matching against %v rules", len(rules))))
	matches := 0
	for _, r := range rules {
		if r.Disabled {
			slog.Debug(e.logTag(tid + fmt.Sprintf("Rule #%v is disabled, skipping", r.Id)))
			continue
		}
		if e.MatchesCondition(tuple, r.Condition, r, tid) {
			slog.Debug(e.logTag(tid+fmt.Sprintf("Rule #%v matches", r.Id)), "name", r.Name)
			matches++
			if r.Throttle.Value == 0 {
				e.ExecuteActions([]types.Message{m}, r, tid)
			} else {
				key := queuesContainer.MakeKey(m.DeviceClass, m.DeviceId, r.Id)
				if !queuesContainer.HasQueue(key) {
					queuesContainer.CreateQueue(key, r.Throttle.Value, func(mm []types.Message) {
						slog.Debug(e.logTag(tid + fmt.Sprintf("Rule #%v messages queue is flushed now", r.Id)))
						e.ExecuteActions(mm, r, tid)
					})
				}
				queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(e.logTag(tid + fmt.Sprintf("Rule #%v message was queued for %s", r.Id, r.Throttle.Value)))
			}
		}
	}
	if matches == 0 {
		slog.Warn(e.logTag(tid + "No one matching rule found"))
	} else {
		slog.Debug(e.logTag(tid + fmt.Sprintf("%v out of %v rules were matched", matches, len(rules))))
	}
	e.ldmService.Set(ldmKey, m)
}
