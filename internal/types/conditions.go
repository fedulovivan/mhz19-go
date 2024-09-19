package types

import (
	"encoding/json"
	"fmt"
)

type CondFn byte

var _ fmt.Stringer = (*CondFn)(nil)
var _ json.Marshaler = (*CondFn)(nil)
var _ json.Unmarshaler = (*CondFn)(nil)

const (
	COND_CHANGED         CondFn = 1
	COND_EQUAL           CondFn = 2
	COND_IN_LIST         CondFn = 3
	COND_IS_NIL          CondFn = 5
	COND_ZIGBEE_DEVICE   CondFn = 6
	COND_DEVICE_CLASS    CondFn = 7
	COND_СHANNEL         CondFn = 8
	COND_FROM_END_DEVICE CondFn = 9
)

var CONDITION_NAMES = map[CondFn]string{
	COND_CHANGED:         "Changed",
	COND_EQUAL:           "Equal",
	COND_IN_LIST:         "InList",
	COND_IS_NIL:          "IsNil",
	COND_ZIGBEE_DEVICE:   "ZigbeeDevice",
	COND_DEVICE_CLASS:    "DeviceClass",
	COND_СHANNEL:         "Channel",
	COND_FROM_END_DEVICE: "FromEndDevice",
}

func (fn CondFn) String() string {
	return CONDITION_NAMES[fn]
}

func (fn CondFn) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, fn)), nil
}

func (fn *CondFn) UnmarshalJSON(b []byte) (err error) {
	var v any
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	for cond, name := range CONDITION_NAMES {
		switch vtyped := v.(type) {
		case string:
			if name == vtyped {
				*fn = cond
				return
			}
		case float64:
			if float64(cond) == vtyped {
				*fn = CondFn(vtyped)
				return
			}
		}
	}
	return fmt.Errorf("failed to unmarshal %v(%T) to CondFn", v, v)
}

type CondImpl func(mt MessageTuple, args Args) (bool, error)

type CondImpls map[CondFn]CondImpl
