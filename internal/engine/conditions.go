package engine

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type CondFn byte

type CondImpl func(mt types.MessageTuple, args Args, e *engine) bool

type CondImpls map[CondFn]CondImpl

const (
	COND_UNKNOWN       CondFn = 0
	COND_CHANGED       CondFn = 1
	COND_EQUAL         CondFn = 2
	COND_IN_LIST       CondFn = 3
	COND_NOT_EQUAL     CondFn = 4
	COND_NOT_NIL       CondFn = 5
	COND_ZIGBEE_DEVICE CondFn = 6
)

var CONDITION_NAMES = map[CondFn]string{
	COND_UNKNOWN:       "<unknown>",
	COND_CHANGED:       "Changed",
	COND_EQUAL:         "Equal",
	COND_IN_LIST:       "InList",
	COND_NOT_EQUAL:     "NotEqual",
	COND_NOT_NIL:       "NotNil",
	COND_ZIGBEE_DEVICE: "ZigbeeDevice",
}

func (s CondFn) String() string {
	return fmt.Sprintf("%v (id=%d)", CONDITION_NAMES[s], s)
}

func (s *CondFn) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, CONDITION_NAMES[*s])), nil
}

var Equal CondImpl = func(mt types.MessageTuple, args Args, e *engine) bool {
	c := NewArgReader(mt[0], args, nil)
	left := c.Get("Left")
	right := c.Get("Right")
	if c.Ok() {
		return left == right
	}
	slog.Error(fmt.Sprintf("Equal: %v", c.Error()))
	return false
}

var NotEqual CondImpl = func(mt types.MessageTuple, args Args, e *engine) bool {
	c := NewArgReader(mt[0], args, nil)
	left := c.Get("Left")
	right := c.Get("Right")
	if c.Ok() {
		return left != right
	}
	slog.Error(fmt.Sprintf("NotEqual: %v", c.Error()))
	return false
}

var InList CondImpl = func(mt types.MessageTuple, args Args, e *engine) bool {
	c := NewArgReader(mt[0], args, nil)
	v := c.Get("Value")
	list := c.Get("List")
	if !c.Ok() {
		slog.Error(fmt.Sprintf("InList: %v", c.Error()))
		return false
	}
	lslice, ok := list.([]any)
	if !ok {
		panic("[]any is expected")
	}
	res := slices.Contains(lslice, v)
	return res
}

// return false for nil and empty strings
// return true for the rest
var NotNil CondImpl = func(mt types.MessageTuple, args Args, e *engine) bool {
	c := NewArgReader(mt[0], args, nil)
	v := c.Get("Value")
	if !c.Ok() {
		slog.Error(fmt.Sprintf("NotNil: %v", c.Error()))
		return false
	}
	switch vTyped := v.(type) {
	case string:
		return len(vTyped) > 0
	case nil:
		return false
	default:
		return true
	}
}

var Changed CondImpl = func(mt types.MessageTuple, args Args, e *engine) bool {
	leftr := NewArgReader(mt[0], args, nil)
	rightr := NewArgReader(mt[1], args, nil)
	left := leftr.Get("Value")
	right := rightr.Get("Value")
	if leftr.Ok() && rightr.Ok() {
		return left != right
	}
	slog.Error(fmt.Sprintf("Changed: %v, %v", leftr.Error(), rightr.Error()))
	return false
}

// args: List
var ZigbeeDevice CondImpl = func(mt types.MessageTuple, args Args, e *engine) bool {
	return Equal(
		mt,
		Args{
			"Left":  "$deviceClass",
			"Right": types.DEVICE_CLASS_ZIGBEE_DEVICE,
		},
		e,
	) && InList(
		mt,
		Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
		e,
	)
}

var conditionImplementations = CondImpls{
	COND_CHANGED:       Changed,
	COND_EQUAL:         Equal,
	COND_IN_LIST:       InList,
	COND_NOT_EQUAL:     NotEqual,
	COND_NOT_NIL:       NotNil,
	COND_ZIGBEE_DEVICE: ZigbeeDevice,
}
