package engine

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine/actions"
	"github.com/fedulovivan/mhz19-go/internal/engine/conditions"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/message_queue"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var queuesContainer = message_queue.NewContainer()

type GetProviderFn func(ch types.ChannelType) types.ChannelProvider

var _ types.Engine = (*engine)(nil)
var _ types.EngineAsSupplier = (*engine)(nil)

type engine struct {
	tag            logger.Tag
	providers      []types.ChannelProvider
	rules          []types.Rule
	messageService types.MessagesService
	devicesService types.DevicesService
	ldmService     types.LdmService
}

func NewEngine() types.Engine {
	return &engine{
		tag: logger.NewTag("default"),
	}
}

func (e *engine) SetLogTag(tag logger.Tag) {
	e.tag = tag
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
	slog.Debug(e.tag.F("AppendRules"), "len", len(rules))
	e.rules = append(e.rules, rules...)
}

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
			for message := range provider.Messages() {
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
	panic(fmt.Sprintf("%v provider is not found", ct))
}

func (e *engine) Stop() {
	for _, s := range e.providers {
		s.Stop()
	}
}

func (e *engine) InvokeConditionFunc(mt types.MessageTuple, fn types.CondFn, not bool, args types.Args, tag logger.Tag) bool {
	impl := conditions.Get(fn)
	res, err := impl(mt, args)
	if not {
		res = !res
	}
	slog.Debug(tag.F("condition exec"), "fn", fn, "args", args, "res", res)
	if err != nil {
		slog.Error(fmt.Sprintf("%s error: %s", fn, err))
	}
	return res
}

func (e *engine) InvokeActionFunc(mm []types.Message, a types.Action, tag logger.Tag) {
	impl := actions.Get(a.Fn)
	slog.Debug(tag.F("action exec"), "fn", a.Fn, "args", a.Args)
	go func() {
		err := impl(mm, a.Args, a.Mapping, e)
		if err != nil {
			slog.Error(fmt.Sprintf("%s error: %s", a.Fn, err))
		}
	}()
}

func (e *engine) MatchesListSome(mtcb types.MessageTupleFn, cc []types.Condition, tag logger.Tag) bool {
	for _, c := range cc {
		if e.MatchesCondition(mtcb, c, tag) {
			return true
		}
	}
	return false
}

func (e *engine) MatchesListEvery(mtcb types.MessageTupleFn, cc []types.Condition, tag logger.Tag) bool {
	if len(cc) == 0 {
		return false
	}
	for _, c := range cc {
		if !e.MatchesCondition(mtcb, c, tag) {
			return false
		}
	}
	return true
}

func (e *engine) MatchesCondition(mtcb types.MessageTupleFn, c types.Condition, tag logger.Tag) bool {
	withFn := c.Fn != 0
	withList := len(c.List) > 0
	if withFn && !withList {
		return e.InvokeConditionFunc(
			mtcb(c.OtherDeviceId), c.Fn, c.Not, c.Args, tag,
		)
	} else if withList && !withFn {
		if c.Or {
			return e.MatchesListSome(mtcb, c.List, tag)
		} else {
			return e.MatchesListEvery(mtcb, c.List, tag)
		}
	} else {
		return true
	}
}

func (e *engine) ExecuteActions(mm []types.Message, r types.Rule, tag logger.Tag) {
	slog.Debug(tag.F(fmt.Sprintf("going to execute %v actions", len(r.Actions))))
	for _, a := range r.Actions {
		e.InvokeActionFunc(mm, a, tag)
	}
}

// called simultaneusly upon receiving messages from all providers
// should be thread-safe
func (e *engine) HandleMessage(m types.Message, rules []types.Rule) {
	app.StatsSingleton().EngineMessagesReceived.Inc()
	tag := e.tag.WithTid()
	p := m.Payload
	isSystem := m.DeviceClass == types.DEVICE_CLASS_SYSTEM
	isBridge := m.DeviceClass == types.DEVICE_CLASS_ZIGBEE_BRIDGE
	sonoffAnnounce := m.DeviceClass == types.DEVICE_CLASS_SONOFF_ANNOUNCE
	if isBridge {
		p = "<too big to render>"
	}
	slog.Debug(
		tag.F("New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)
	ldmKey := e.ldmService.NewKey(m.DeviceClass, m.DeviceId)
	slog.Debug(tag.F(fmt.Sprintf("Matching against %v rules", len(rules))))
	matches := 0
	for _, r := range rules {
		tag := tag.With("Rule#%v", r.Id)
		if r.Disabled {
			slog.Debug(tag.F("is disabled, skipping"))
			continue
		}
		var mtcb = func(otherDeviceId types.DeviceId) types.MessageTuple {
			tuple := types.MessageTuple{}
			takeOtherDeviceMessage := len(otherDeviceId) > 0
			if takeOtherDeviceMessage {
				slog.Warn(tag.F(fmt.Sprintf("requesting message for otherDeviceId=%v", otherDeviceId)))
				otherLdmKey := e.ldmService.NewKey(m.DeviceClass, otherDeviceId)
				if e.ldmService.Has(otherLdmKey) {
					tmp := e.ldmService.Get(otherLdmKey)
					tuple.Curr = &tmp
				}
			} else {
				tuple.Curr = &m
				if e.ldmService.Has(ldmKey) {
					tmp := e.ldmService.Get(ldmKey)
					tuple.Prev = &tmp
				}
			}
			return tuple
		}
		if e.MatchesCondition(mtcb, r.Condition, tag) {
			if !r.SkipCounter {
				app.StatsSingleton().EngineRulesMatched.Inc()
			}
			slog.Debug(tag.F("matches"), "name", r.Name)
			matches++
			if r.Throttle.Duration == 0 {
				e.ExecuteActions([]types.Message{m}, r, tag)
			} else {
				key := queuesContainer.MakeKey(m.DeviceClass, m.DeviceId, r.Id)
				if !queuesContainer.HasQueue(key) {
					queuesContainer.CreateQueue(key, r.Throttle.Duration, func(mm []types.Message) {
						slog.Debug(tag.F("messages queue is flushed now"))
						e.ExecuteActions(mm, r, tag)
					})
				}
				queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(tag.F(fmt.Sprintf("message was queued for %s", r.Throttle.Duration)))
			}
		}
	}
	if matches == 0 {
		slog.Warn(tag.F("No one matching rule found"))
	} else {
		slog.Debug(tag.F(fmt.Sprintf("%v out of %v rules were matched", matches, len(rules))))
	}
	if !isBridge && !sonoffAnnounce && !isSystem {
		e.ldmService.Set(ldmKey, m)
	}
}
