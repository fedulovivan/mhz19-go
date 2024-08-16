package engine

import (
	"fmt"
	"slices"
)

type CondFn byte

type CondImpl func(mt MessageTuple, args Args) bool

type CondImpls map[CondFn]CondImpl

const (
	COND_CHANGED       CondFn = 1
	COND_EQUAL         CondFn = 2
	COND_IN_LIST       CondFn = 3
	COND_NOT_EQUAL     CondFn = 4
	COND_NOT_NIL       CondFn = 5
	COND_ZIGBEE_DEVICE CondFn = 6
)

var CONDITION_NAMES = map[CondFn]string{
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

var Equal CondImpl = func(mt MessageTuple, args Args) bool {
	c := NewArgReader(COND_EQUAL, mt[0], args)
	left := c.Get("Left")
	right := c.Get("Right")
	// left := Get[int](c, "Left")
	// right := Get[int](c, "Right")
	if c.Ok() {
		return left == right
	}
	return false
}

var NotEqual CondImpl = func(mt MessageTuple, args Args) bool {
	c := NewArgReader(COND_NOT_EQUAL, mt[0], args)
	left := c.Get("Left")
	right := c.Get("Right")
	if c.Ok() {
		return left != right
	}
	return false
}

var InList CondImpl = func(mt MessageTuple, args Args) bool {
	c := NewArgReader(COND_IN_LIST, mt[0], args)
	v := c.Get("Value")
	list := c.Get("List")
	if !c.Ok() {
		return false
	}
	lslice, ok := list.([]any)
	if !ok {
		panic("[]any is expected")
	}
	res := slices.Contains(lslice, v)
	// fmt.Printf("InList: %v: %T in %T, %v and %v\n", res, v, lslice, v, lslice)
	return res
}

// return false for nil and empty strings
// return true for the rest
var NotNil CondImpl = func(mt MessageTuple, args Args) bool {
	c := NewArgReader(COND_NOT_NIL, mt[0], args)
	v := c.Get("Value")
	if !c.Ok() {
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

var Changed CondImpl = func(mt MessageTuple, args Args) bool {
	leftr := NewArgReader(COND_CHANGED, mt[0], args)
	rightr := NewArgReader(COND_CHANGED, mt[1], args)
	left := leftr.Get("Value")
	right := rightr.Get("Value")
	if leftr.Ok() && rightr.Ok() {
		return left != right
	}
	return false
}

var ZigbeeDevice CondImpl = func(mt MessageTuple, args Args) bool {
	return Equal(
		mt,
		Args{
			"Left":  "$deviceClass",
			"Right": DEVICE_CLASS_ZIGBEE_DEVICE,
		},
	) && InList(
		mt,
		Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
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

// func GEqual[T comparable](a T, b T) bool {
// 	return a == b
// }
// func GInList[T comparable](list []T, v T) bool {
// 	return slices.Contains(list, v)
// }
// func GNotNil(v any) bool {
// 	switch vTyped := v.(type) {
// 	case string:
// 		return len(vTyped) > 0
// 	case nil:
// 		return false
// 	default:
// 		return true
// 	}
// }
