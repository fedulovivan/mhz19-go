package engine

import (
	"time"
)

type JsonPayload any

type DeviceId string

type Service interface {
	Receive() MessageChan
	// Type() ChannelType
	Init()
	Stop()
}

type NamedArgs map[string]any

type Rule struct {
	Condition Condition
	Actions   []Action
	Throttle  time.Duration
}

type Condition struct {
	Fn   CondFnName
	Args NamedArgs
	List []Condition
	Or   bool
}

type Action struct {
	Fn   ActionFnName
	Args NamedArgs
	Mapping
}

// {
// 	Fn: ZigbeeDevice,
// 	Args: NamedArgs{
// 		"DeviceId": []string{"0x00158d0004244bda"},
// 	},
// },
