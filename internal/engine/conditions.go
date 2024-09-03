package engine

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Equal types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	c := arg_reader.NewArgReader(mt[0], args, nil, nil)
	left := c.Get("Left")
	right := c.Get("Right")
	if c.Ok() {
		return left == right
	}
	slog.Error(fmt.Sprintf("Equal: %v", c.Error()))
	return false
}

var NotEqual types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	c := arg_reader.NewArgReader(mt[0], args, nil, nil)
	left := c.Get("Left")
	right := c.Get("Right")
	if c.Ok() {
		return left != right
	}
	slog.Error(fmt.Sprintf("NotEqual: %v", c.Error()))
	return false
}

var InList types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	c := arg_reader.NewArgReader(mt[0], args, nil, nil)
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
var NotNil types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	c := arg_reader.NewArgReader(mt[0], args, nil, nil)
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

var Changed types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	leftr := arg_reader.NewArgReader(mt[0], args, nil, nil)
	rightr := arg_reader.NewArgReader(mt[1], args, nil, nil)
	left := leftr.Get("Value")
	right := rightr.Get("Value")
	if leftr.Ok() && rightr.Ok() {
		return left != right
	}
	slog.Error(fmt.Sprintf("Changed: %v, %v", leftr.Error(), rightr.Error()))
	return false
}

// args: List
var ZigbeeDevice types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	return Equal(
		mt,
		types.Args{
			"Left":  "$deviceClass",
			"Right": types.DEVICE_CLASS_ZIGBEE_DEVICE,
		},
		e,
	) && InList(
		mt,
		types.Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
		e,
	)
}

var conditionImplementations = types.CondImpls{
	types.COND_CHANGED:       Changed,
	types.COND_EQUAL:         Equal,
	types.COND_IN_LIST:       InList,
	types.COND_NOT_EQUAL:     NotEqual,
	types.COND_NOT_NIL:       NotNil,
	types.COND_ZIGBEE_DEVICE: ZigbeeDevice,
}
