package engine

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/message_queue"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var queuesContainer = message_queue.NewContainer()

type ChannelProvider interface {

	// getter for the provider's messages channel used in engine.Start
	MessageChan() types.MessageChan

	// api to invoke provider outbound action, eg:
	// - call tgbotapi.NewMessage for telegram bot provider
	// - post to mqtt topic for mqtt provider
	// - call sonoff http api
	Send(...any)

	// TODO, api for the unit tests
	Write(m types.Message)

	Channel() types.ChannelType
	Init()
	Stop()
}

var tidSeq = utils.NewSeq()

type engine struct {
	options Options
}

func NewEngine(
	options Options,
) engine {
	return engine{
		options: options,
	}
}

func (e engine) Start() {
	for _, p := range e.options.providers {
		go func(provider ChannelProvider) {
			provider.Init()
			for mchan := range provider.MessageChan() {
				e.handleMessage(mchan, e.options.rules)
			}
		}(p)
	}
}

func (e engine) getPrivider(ct types.ChannelType) ChannelProvider {
	for _, provider := range e.options.providers {
		if provider.Channel() == ct {
			return provider
		}
	}
	return nil
}

func (e engine) Stop() {
	for _, s := range e.options.providers {
		s.Stop()
	}
}

func (e engine) invokeConditionFunc(mt types.MessageTuple, fn CondFn, args Args, r Rule, tid string) bool {
	impl, ok := conditionImplementations[fn]
	if !ok {
		slog.Error(e.options.logTag(fmt.Sprintf("Condition function [%v] not yet implemented", fn)))
		return false
	}
	res := impl(mt, args, &e)
	slog.Debug(e.options.logTag(tid+fmt.Sprintf("Rule #%v condition exec", r.Id)), "fn", fn, "args", args, "res", res)
	return res
}

func (e engine) invokeActionFunc(mm []types.Message, a Action, r Rule, tid string) {
	impl, ok := actions[a.Fn]
	if !ok {
		slog.Error(e.options.logTag(fmt.Sprintf("Action function [%v] not yet implemented", a.Fn)))
		return
	}
	slog.Debug(e.options.logTag(tid+fmt.Sprintf("Rule #%v action exec", r.Id)), "fn", a.Fn, "args", a.Args)
	go impl(mm, a /* e.getPrivider,  */, &e)
}

func (e engine) matchesListSome(mt types.MessageTuple, cc []Condition, r Rule, tid string) bool {
	for _, c := range cc {
		if e.matchesCondition(mt, c, r, tid) {
			return true
		}
	}
	return false
}

func (e engine) matchesListEvery(mt types.MessageTuple, cc []Condition, r Rule, tid string) bool {
	if len(cc) == 0 {
		return false
	}
	for _, c := range cc {
		if !e.matchesCondition(mt, c, r, tid) {
			return false
		}
	}
	return true
}

func (e engine) matchesCondition(mt types.MessageTuple, c Condition, r Rule, tid string) bool {
	withFn := c.Fn != 0
	withList := len(c.List) > 0
	if withFn && !withList {
		return e.invokeConditionFunc(mt, c.Fn, c.Args, r, tid)
	} else if withList && !withFn {
		if c.Or {
			return e.matchesListSome(mt, c.List, r, tid)
		} else {
			return e.matchesListEvery(mt, c.List, r, tid)
		}
	} else {
		return true
	}
}

func (e engine) executeActions(mm []types.Message, r Rule, tid string) {
	slog.Debug(e.options.logTag(tid + fmt.Sprintf("Rule #%v going to execute %v actions", r.Id, len(r.Actions))))
	for _, a := range r.Actions {
		e.invokeActionFunc(mm, a, r, tid)
	}
}

// called simultaneusly upon receiving messages from all providers
// should be thread-safe
func (e engine) handleMessage(m types.Message, rules []Rule) {
	tid := fmt.Sprintf("Tid #%v ", tidSeq.Next())
	p := m.Payload
	if m.DeviceClass == types.DEVICE_CLASS_ZIGBEE_BRIDGE {
		p = "<too big to render>"
	}
	slog.Debug(
		e.options.logTag(tid+"New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)
	mkey := makeMessageKey(m.DeviceClass, m.DeviceId)
	prevm := PrevMessageGet(mkey)
	mt := types.MessageTuple{m, prevm}
	slog.Debug(e.options.logTag(tid + fmt.Sprintf("Matching against %v rules", len(rules))))
	matches := 0
	for _, r := range rules {
		if r.Disabled {
			slog.Debug(e.options.logTag(tid + fmt.Sprintf("Rule #%v is disabled, skipping", r.Id)))
			continue
		}
		if e.matchesCondition(mt, r.Condition, r, tid) {
			slog.Debug(e.options.logTag(tid+fmt.Sprintf("Rule #%v matches", r.Id)), "comments", r.Comments)
			matches++
			if r.Throttle == 0 {
				e.executeActions(mt[:], r, tid)
			} else {
				key := queuesContainer.MakeKey(m.DeviceClass, m.DeviceId, r.Id)
				if !queuesContainer.HasQueue(key) {
					queuesContainer.CreateQueue(key, r.Throttle, func(mm []types.Message) {
						slog.Debug(e.options.logTag(tid + fmt.Sprintf("Rule #%v messages queue is flushed now", r.Id)))
						e.executeActions(mm, r, tid)
					})
				}
				queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(e.options.logTag(tid + fmt.Sprintf("Rule #%v message was queued", r.Id)))
			}
		}
	}
	if matches == 0 {
		slog.Warn(e.options.logTag(tid + "No one matching rule found"))
	} else {
		slog.Debug(e.options.logTag(tid + fmt.Sprintf("%v out of %v rules were matched", matches, len(rules))))
	}
	if prevMessages == nil {
		prevMessages = make(map[string]types.Message)
	}
	PrevMessagePut(mkey, m)
}
