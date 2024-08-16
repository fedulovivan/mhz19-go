package engine

import (
	"fmt"
	"log/slog"
	"sync"
)

type autoInc struct {
	sync.Mutex
	id int
}

func (a *autoInc) Next() (id int) {
	a.Lock()
	defer a.Unlock()
	a.id++
	id = a.id
	return
}

var ai autoInc

func Start(o Options) {
	opts = o
	start()
}

func start() {
	for _, service := range opts.services {
		go func(s Service) {
			s.Init()
			for m := range s.Receive() {
				handleMessage(m, Rules)
			}
		}(service)
	}
}

func getService(ct ChannelType) Service {
	for _, service := range opts.services {
		if service.Channel() == ct {
			return service
		}
	}
	return nil
}

func Stop() {
	for _, s := range opts.services {
		s.Stop()
	}
}

func invokeConditionFunc(mt MessageTuple, fn CondFn, args Args, r Rule, tid string) bool {
	impl, ok := conditionImplementations[fn]
	if !ok {
		slog.Error(fmt.Sprintf("Condition function [%v] not yet implemented", fn))
		return false
	}
	res := impl(mt, args)
	slog.Debug(opts.logTag(tid+fmt.Sprintf("Rule #%v condition exec", r.Id)), "fn", fn, "args", args, "res", res)
	return res
}

func invokeActionFunc(mm []Message, a Action, r Rule, tid string) {
	impl, ok := actions[a.Fn]
	if !ok {
		slog.Error(fmt.Sprintf("Action function [%v] not yet implemented", a.Fn))
		return
	}
	slog.Debug(opts.logTag(tid+fmt.Sprintf("Rule #%v action exec", r.Id)), "fn", a.Fn, "args", a.Args)
	go impl(mm, a, getService)
}

func matchesListSome(mt MessageTuple, cc []Condition, r Rule, tid string) bool {
	for _, c := range cc {
		if matchesCondition(mt, c, r, tid) {
			return true
		}
	}
	return false
}

func matchesListEvery(mt MessageTuple, cc []Condition, r Rule, tid string) bool {
	if len(cc) == 0 {
		return false
	}
	for _, c := range cc {
		if !matchesCondition(mt, c, r, tid) {
			return false
		}
	}
	return true
}

func matchesCondition(mt MessageTuple, c Condition, r Rule, tid string) bool {
	withFn := c.Fn != 0
	withList := len(c.List) > 0
	if withFn && !withList {
		return invokeConditionFunc(mt, c.Fn, c.Args, r, tid)
	} else if withList && !withFn {
		if c.Or {
			return matchesListSome(mt, c.List, r, tid)
		} else {
			return matchesListEvery(mt, c.List, r, tid)
		}
	} else {
		return false
	}
}

func executeActions(mm []Message, aa []Action, r Rule, tid string) {
	slog.Debug(opts.logTag(tid + fmt.Sprintf("Rule #%v going to execute %v actions", r.Id, len(aa))))
	for _, a := range aa {
		invokeActionFunc(mm, a, r, tid)
	}
}

func handleMessage(m Message, rules []Rule) {
	tid := fmt.Sprintf("Tid #%v ", ai.Next())
	p := m.Payload
	if m.DeviceClass == DEVICE_CLASS_ZIGBEE_BRIDGE {
		p = "<too big to render>"
	}
	slog.Debug(
		opts.logTag(tid+"New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)
	mkey := makeMessageKey(m.DeviceClass, m.DeviceId)
	prevm := PrevMessageGet(mkey)
	mt := MessageTuple{m, prevm}
	slog.Debug(opts.logTag(tid + fmt.Sprintf("Matching against %v rules", len(rules))))
	matches := 0
	for _, r := range rules {
		if r.Disabled {
			slog.Debug(opts.logTag(tid + fmt.Sprintf("Rule #%v is disabled, skipping", r.Id)))
			continue
		}
		if matchesCondition(mt, r.Condition, r, tid) {
			slog.Debug(opts.logTag(tid+fmt.Sprintf("Rule #%v matches", r.Id)), "comments", r.Comments)
			matches++
			if r.Throttle == 0 {
				executeActions(mt[:], r.Actions, r, tid)
			} else {
				panic("handleMessage: not yet implemented: Throttle")
			}
		}
	}
	if matches == 0 {
		slog.Warn(opts.logTag(tid + "No one matching rule found"))
	} else {
		slog.Debug(opts.logTag(tid + fmt.Sprintf("%v out of %v rules were matched", matches, len(rules))))
	}
	if prevMessages == nil {
		prevMessages = make(map[string]Message)
	}
	PrevMessagePut(mkey, m)
}
