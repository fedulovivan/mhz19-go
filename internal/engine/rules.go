package engine

import (
	"encoding/json"
	"strconv"
	"strings"
	"time"
)

type DeviceId string

// func (d DeviceId) MarshalJSON() ([]byte, error) {
// 	return []byte(fmt.Sprintf(`"DeviceId(%s)"`, d)), nil
// }

type Args map[string]any

func parseDeviceIdOrClass(in string) any {
	if strings.HasPrefix(in, "DeviceId(") {
		deviceId := in[9 : len(in)-1]
		return DeviceId(deviceId)
	} else if strings.HasPrefix(in, "DeviceClass(") {
		deviceClass := in[12 : len(in)-1]
		i, _ := strconv.Atoi(deviceClass)
		return DeviceClass(i)
	}
	return in
}

// TODO seems there is a room for optimization here
func (a *Args) UnmarshalJSON(data []byte) (err error) {
	var raw map[string]any
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	*a = make(map[string]any, len(raw))
	for argName, argValue := range raw {
		switch vtyped := argValue.(type) {
		case []any:
			for i, listel := range vtyped {
				if slistel, ok := listel.(string); ok {
					vtyped[i] = parseDeviceIdOrClass(slistel)
				}
			}
			(*a)[argName] = vtyped
		case string:
			(*a)[argName] = parseDeviceIdOrClass(vtyped)
		default:
			(*a)[argName] = vtyped
		}
	}
	return
}

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
	Disabled  bool          `json:"disabled,omitempty"`
	Comments  string        `json:"comments,omitempty"`
	Condition Condition     `json:"condition,omitempty"`
	Actions   []Action      `json:"actions,omitempty"`
	Throttle  time.Duration `json:"throttle,omitempty"`
}

type Condition struct {
	Id   int         `json:"-"`
	Fn   CondFn      `json:"fn,omitempty"`
	Args Args        `json:"args,omitempty"`
	List []Condition `json:"list,omitempty"`
	Or   bool        `json:"or,omitempty"`
}

type Action struct {
	Id       int      `json:"-"`
	Fn       ActionFn `json:"fn,omitempty"`
	Args     Args     `json:"args,omitempty"`
	Mapping  Mapping  `json:"mapping,omitempty"`
	DeviceId DeviceId `json:"deviceId,omitempty"`
}
