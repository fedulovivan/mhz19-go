package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

type Args map[string]any

func Value(value any) Args {
	return Args{
		"Value": value,
	}
}

func parseDeviceIdOrClass(in string) any {
	if strings.HasPrefix(in, "DeviceId(") {
		deviceId := in[9 : len(in)-1]
		return DeviceId(deviceId)
	} else if strings.HasPrefix(in, "DeviceClass(") {
		dc := in[12 : len(in)-1]
		i, _ := strconv.Atoi(dc)
		return DeviceClass(i)
	} else if strings.HasPrefix(in, "ChannelType(") {
		ct := in[12 : len(in)-1]
		i, _ := strconv.Atoi(ct)
		return ChannelType(i)
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

type TemplatePayload struct {
	Message  Message
	Messages []Message
}
