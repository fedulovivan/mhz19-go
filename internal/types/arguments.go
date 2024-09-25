package types

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"
)

type Args map[string]any

func Value(value any) Args {
	return Args{
		"Value": value,
	}
}

func parseSpecial(in string) (res any, err error) {

	if strings.HasPrefix(in, "DeviceId(") {
		var typed DeviceId
		err = json.Unmarshal(
			[]byte(fmt.Sprintf(`"%s"`, in)),
			&typed,
		)
		if err == nil {
			res = typed
		}
		return
	}

	if strings.HasPrefix(in, "DeviceClass(") {
		dc := in[12 : len(in)-1]
		i, _ := strconv.Atoi(dc)
		res = DeviceClass(i)
		return
	}

	if strings.HasPrefix(in, "ChannelType(") {
		ct := in[12 : len(in)-1]
		i, _ := strconv.Atoi(ct)
		res = ChannelType(i)
		return
	}

	if strings.HasPrefix(in, "Channel(") { // same as ChannelType
		ct := in[8 : len(in)-1]
		i, _ := strconv.Atoi(ct)
		res = ChannelType(i)
		return
	}

	res = in
	return
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
					vtyped[i], err = parseSpecial(slistel)
				}
			}
			(*a)[argName] = vtyped
		case string:
			(*a)[argName], err = parseSpecial(vtyped)
		default:
			(*a)[argName] = vtyped
		}
	}
	return
}

type TemplatePayload struct {
	IsFirst  bool
	Message  Message
	Messages []Message
}
