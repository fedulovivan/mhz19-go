package types

type Mapping map[string](map[string]string)

type Rule struct {
	Id        int       `json:"id"`
	Disabled  bool      `json:"disabled,omitempty"`
	Name      string    `json:"name,omitempty"`
	Condition Condition `json:"condition,omitempty"`
	Actions   []Action  `json:"actions,omitempty"`
	Throttle  Throttle  `json:"throttle,omitempty"`
}

type Condition struct {
	Id            int         `json:"-"`
	Fn            CondFn      `json:"fn,omitempty"`
	Args          Args        `json:"args,omitempty"`
	Nested        []Condition `json:"nested,omitempty"`
	Or            bool        `json:"or,omitempty"`
	Not           bool        `json:"not,omitempty"`
	OtherDeviceId DeviceId    `json:"otherDeviceId,omitempty"`
}

type Action struct {
	Id      int      `json:"-"`
	Fn      ActionFn `json:"fn,omitempty"`
	Args    Args     `json:"args,omitempty"`
	Mapping Mapping  `json:"mapping,omitempty"`
}
