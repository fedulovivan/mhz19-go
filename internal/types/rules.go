package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Mapping map[string](map[string]string)

type Throttle struct {
	Value time.Duration
}

func (t *Throttle) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, t.Value)), nil
}

func (t *Throttle) UnmarshalJSON(b []byte) (err error) {
	var v any
	err = json.Unmarshal(b, &v)
	if err != nil {
		return
	}
	vstring, issting := v.(string)
	if !issting {
		return
	}
	t.Value, err = time.ParseDuration(vstring)
	if err != nil {
		return
	}
	return
}

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
	List          []Condition `json:"list,omitempty"`
	Or            bool        `json:"or,omitempty"`
	OtherDeviceId DeviceId    `json:"otherDeviceId,omitempty"`
}

type Action struct {
	Id      int      `json:"-"`
	Fn      ActionFn `json:"fn,omitempty"`
	Args    Args     `json:"args,omitempty"`
	Mapping Mapping  `json:"mapping,omitempty"`
}
