package engine

import (
	"time"
)

type DeviceId string

type Args map[string]any

type Mapping map[string](map[string]string)

type Provider interface {
	Receive() MessageChan
	Send(...any)
	Channel() ChannelType
	Init()
	Stop()
}

type Rule struct {
	Id        int32         `json:"id"`
	Disabled  bool          `json:"disabled"`
	Comments  string        `json:"comments"`
	Condition Condition     `json:"condition"`
	Actions   []Action      `json:"actions"`
	Throttle  time.Duration `json:"throttle"`
}

type Condition struct {
	Id   int         `json:"-"`
	Fn   CondFn      `json:"fn,omitempty"`
	Args Args        `json:"args,omitempty"`
	List []Condition `json:"list,omitempty"`
	Or   bool        `json:"or,omitempty"`
}

type Action struct {
	Id      int      `json:"-"`
	Fn      ActionFn `json:"fn"`
	Args    Args     `json:"args"`
	Mapping Mapping  `json:"mapping"`
}
