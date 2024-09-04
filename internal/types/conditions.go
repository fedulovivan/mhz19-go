package types

import "fmt"

type CondFn byte

const (
	COND_UNKNOWN       CondFn = 0
	COND_CHANGED       CondFn = 1
	COND_EQUAL         CondFn = 2
	COND_IN_LIST       CondFn = 3
	COND_NOT_EQUAL     CondFn = 4
	COND_NOT_NIL       CondFn = 5
	COND_ZIGBEE_DEVICE CondFn = 6
	COND_DEVICE_CLASS  CondFn = 7
)

var CONDITION_NAMES = map[CondFn]string{
	COND_UNKNOWN:       "<unknown>",
	COND_CHANGED:       "Changed",
	COND_EQUAL:         "Equal",
	COND_IN_LIST:       "InList",
	COND_NOT_EQUAL:     "NotEqual",
	COND_NOT_NIL:       "NotNil",
	COND_ZIGBEE_DEVICE: "ZigbeeDevice",
	COND_DEVICE_CLASS:  "DeviceClass",
}

func (s CondFn) String() string {
	return fmt.Sprintf("%v (id=%d)", CONDITION_NAMES[s], s)
}

func (s *CondFn) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, CONDITION_NAMES[*s])), nil
}

type CondImpl func(mt MessageTuple, args Args, e Engine) bool

type CondImpls map[CondFn]CondImpl
