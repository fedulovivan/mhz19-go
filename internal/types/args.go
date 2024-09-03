package types

import (
	"encoding/json"
	"strconv"
	"strings"
)

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
