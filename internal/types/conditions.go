package types

import (
	"encoding/json"
	"fmt"

	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type CondFn byte

var _ fmt.Stringer = (*CondFn)(nil)
var _ json.Marshaler = (*CondFn)(nil)
var _ json.Unmarshaler = (*CondFn)(nil)

const (
	COND_CHANGED         CondFn = 1
	COND_EQUAL           CondFn = 2
	COND_IN_LIST         CondFn = 3
	COND_NIL             CondFn = 5
	COND_ZIGBEE_DEVICE   CondFn = 6
	COND_DEVICE_CLASS    CondFn = 7
	COND_СHANNEL         CondFn = 8
	COND_FROM_END_DEVICE CondFn = 9
	COND_TRUE            CondFn = 10
	COND_FALSE           CondFn = 11
	COND_DEVICE_ID       CondFn = 12
	COND_LDM_OLDER_THAN  CondFn = 13
)

var CONDITION_NAMES = map[CondFn]string{
	COND_CHANGED:         "Changed",
	COND_EQUAL:           "Equal",
	COND_IN_LIST:         "InList",
	COND_NIL:             "Nil",
	COND_ZIGBEE_DEVICE:   "ZigbeeDevice",
	COND_DEVICE_CLASS:    "DeviceClass",
	COND_СHANNEL:         "Channel",
	COND_FROM_END_DEVICE: "FromEndDevice",
	COND_TRUE:            "True",
	COND_FALSE:           "False",
	COND_DEVICE_ID:       "DeviceId",
	COND_LDM_OLDER_THAN:  "LdmOlderThan",
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
	return fmt.Errorf("failed to unmarshal %v (type=%T) to CondFn", v, v)
}

type CondImpl func(mt MessageCompound, args Args, tag utils.Tag) (bool, error)

type CondImpls map[CondFn]CondImpl
