package engine

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/engine_actions"
	"github.com/fedulovivan/mhz19-go/internal/message_queue"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var queuesContainer = message_queue.NewContainer()

type GetProviderFn func(ch types.ChannelType) types.ChannelProvider

var tidSeq = utils.NewSeq()

type engine struct {
	options types.EngineOptions
}

func NewEngine(
	options types.EngineOptions,
) types.Engine {
	return &engine{
		options: options,
	}
}

func (e *engine) GetOptions() types.EngineOptions {
	return e.options
}

func (e *engine) Start() {
	for _, p := range e.options.Providers() {
		go func(provider types.ChannelProvider) {
			provider.Init()
			for mchan := range provider.MessageChan() {
				e.HandleMessage(mchan, e.options.Rules())
			}
		}(p)
	}
}

func (e *engine) FindProvider(ct types.ChannelType) types.ChannelProvider {
	for _, provider := range e.options.Providers() {
		if provider.Channel() == ct {
			return provider
		}
	}
	return nil
}

func (e *engine) Stop() {
	for _, s := range e.options.Providers() {
		s.Stop()
	}
}

func (e *engine) InvokeConditionFunc(mt types.MessageTuple, fn types.CondFn, args types.Args, r types.Rule, tid string) bool {
	impl, ok := conditionImplementations[fn]
	logTag := e.GetOptions().LogTag()
	if !ok {
		slog.Error(logTag(fmt.Sprintf("Condition function [%v] not yet implemented", fn)))
		return false
	}
	res := impl(mt, args, e)
	slog.Debug(logTag(tid+fmt.Sprintf("Rule #%v condition exec", r.Id)), "fn", fn, "args", args, "res", res)
	return res
}

func (e *engine) invokeActionFunc(mm []types.Message, a types.Action, r types.Rule, tid string) {
	impl, ok := engine_actions.Actions[a.Fn]
	logTag := e.GetOptions().LogTag()
	if !ok {
		slog.Error(logTag(fmt.Sprintf("Action function [%v] not yet implemented", a.Fn)))
		return
	}
	slog.Debug(logTag(tid+fmt.Sprintf("Rule #%v action exec", r.Id)), "fn", a.Fn, "args", a.Args)
	go impl(mm, a, e)
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
	logTag := e.GetOptions().LogTag()
	slog.Debug(logTag(tid + fmt.Sprintf("Rule #%v going to execute %v actions", r.Id, len(r.Actions))))
	for _, a := range r.Actions {
		e.invokeActionFunc(mm, a, r, tid)
	}
}

// called simultaneusly upon receiving messages from all providers
// should be thread-safe
func (e *engine) HandleMessage(m types.Message, rules []types.Rule) {
	logTag := e.GetOptions().LogTag()
	ldms := e.GetOptions().LdmService()
	tid := fmt.Sprintf("Tid #%v ", tidSeq.Next())
	p := m.Payload
	if m.DeviceClass == types.DEVICE_CLASS_ZIGBEE_BRIDGE {
		p = "<too big to render>"
	}
	slog.Debug(
		logTag(tid+"New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)
	ldmKey := ldms.MakeKey(m.DeviceClass, m.DeviceId)
	prevmessage := ldms.Get(ldmKey)
	mt := types.MessageTuple{m, prevmessage}
	slog.Debug(logTag(tid + fmt.Sprintf("Matching against %v rules", len(rules))))
	matches := 0
	for _, r := range rules {
		if r.Disabled {
			slog.Debug(logTag(tid + fmt.Sprintf("Rule #%v is disabled, skipping", r.Id)))
			continue
		}
		if e.MatchesCondition(mt, r.Condition, r, tid) {
			slog.Debug(logTag(tid+fmt.Sprintf("Rule #%v matches", r.Id)), "name", r.Name)
			matches++
			if r.Throttle == 0 {
				e.ExecuteActions(mt[:], r, tid)
			} else {
				key := queuesContainer.MakeKey(m.DeviceClass, m.DeviceId, r.Id)
				if !queuesContainer.HasQueue(key) {
					queuesContainer.CreateQueue(key, r.Throttle, func(mm []types.Message) {
						slog.Debug(logTag(tid + fmt.Sprintf("Rule #%v messages queue is flushed now", r.Id)))
						e.ExecuteActions(mm, r, tid)
					})
				}
				queuesContainer.GetQueue(key).PushMessage(m)
				slog.Debug(logTag(tid + fmt.Sprintf("Rule #%v message was queued", r.Id)))
			}
		}
	}
	if matches == 0 {
		slog.Warn(logTag(tid + "No one matching rule found"))
	} else {
		slog.Debug(logTag(tid + fmt.Sprintf("%v out of %v rules were matched", matches, len(rules))))
	}
	ldms.Set(ldmKey, m)
}
