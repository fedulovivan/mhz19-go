package engine

import (
	"time"
)

type JsonPayload any

type DeviceId string

type Service interface {
	Receive() MessageChan
	Send(...any)
	Channel() ChannelType
	Init()
	Stop()
}

type Rule struct {
	Id        int
	Disabled  bool
	Comments  string
	Condition Condition
	Actions   []Action
	Throttle  time.Duration
}

type Condition struct {
	Id   int
	Fn   CondFn
	Args Args
	List []Condition
	Or   bool
}

type Action struct {
	Id      int
	Fn      ActionFn
	Args    Args
	Mapping Mapping
}
