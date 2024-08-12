package engine

import (
	"fmt"
	"log/slog"
	"strings"

	"github.com/fedulovivan/mhz19-go/internal/logger"
)

var withTag = logger.MakeTag("ENGN")

var services []Service

func Start(input ...Service) {
	services = input
	start()
}

func start() {
	for _, service := range services {
		go func(s Service) {
			s.Init()
			for m := range s.Receive() {
				handleMessage(m, Rules)
			}
		}(service)
	}
}

func Stop() {
	for _, s := range services {
		s.Stop()
	}
}

func arg(m Message, args NamedArgs, argName string) (any, error) {
	// check such arg exist
	v, ok := args[argName]
	if !ok {
		return nil, fmt.Errorf("no such argument %v", argName)
	}
	// check argument is string
	stringed, ok := v.(string)
	if !ok {
		return v, nil
	}
	// parse value with directive
	if strings.HasPrefix(stringed, "$message.") {
		_, field, _ := strings.Cut(stringed, ".")
		return m.Get(field)
		// if ok {
		// }
		// return stringed, nil
	} else if stringed == "$deviceId" {
		return m.DeviceId, nil
	} else if stringed == "$deviceClass" {
		return m.DeviceClass, nil
	}
	return stringed, nil
}

func matchFunction(mt MessageTuple, fn CondFnName, args NamedArgs) bool {
	impl, ok := conditions[fn]
	if !ok {
		panic(fmt.Sprintf("matchFunction: not yet implemented: %v", fn))
	}
	return impl(mt, args)
}

func listSome(mt MessageTuple, cc []Condition) bool {
	for _, c := range cc {
		if matchesCondition(mt, c) {
			return true
		}
	}
	return false
}

func listEvery(mt MessageTuple, cc []Condition) bool {
	matches := 0
	for _, c := range cc {
		if matchesCondition(mt, c) {
			matches++
		}
	}
	return matches > 0 && len(cc) == matches
}

func matchesCondition(mt MessageTuple, c Condition) bool {
	withFn := c.Fn != ""
	withList := len(c.List) > 0
	if withFn && !withList {
		return matchFunction(mt, c.Fn, c.Args)
	} else if withList && !withFn {
		if c.Or {
			return listSome(mt, c.List)
		} else {
			return listEvery(mt, c.List)
		}
	} else {
		return false
	}
}

func executeActions(a []Action) {
	slog.Debug(withTag("going to execute"), "actions", a)
}

func handleMessage(m Message, rules []Rule) {

	p := m.Payload
	if m.DeviceClass == DEVICE_CLASS_ZIGBEE_BRIDGE {
		p = "<too big to render>"
	}

	slog.Debug(
		withTag("New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)

	key := makeMessageKey(m.DeviceClass, m.DeviceId)

	prevm := PrevMessageGet(key)

	mt := MessageTuple{m, prevm}

	for _, r := range rules {
		if matchesCondition(mt, r.Condition) {
			if r.Throttle == 0 {
				executeActions(r.Actions)
			} else {
				panic("handleMessage: not yet implemented: Throttle")
			}
		}
	}

	if prevMessages == nil {
		prevMessages = make(map[string]Message)
	}

	PrevMessagePut(key, m)

}
