package types

import (
	"encoding/json"
	"strings"
)

// todo replace Args with struct(s) like
//
//	type PostSonoffSwitchMessage_Args struct {
//		DeviceId DeviceId
//		Command  string
//	}
//
// (!) consider counterpart code from internal/entities/rules/service.go::BuildArguments
type Args map[string]any

func (a *Args) UnmarshalJSON(data []byte) error {
	// step 1: parse json into untyped map (json parsing error will be handled earlier - so no need to bother here)
	var raw map[string]any
	_ = json.Unmarshal(data, &raw)
	// step 2: iterate and replace what possible (only two levels are supported, no recursion)
	var err error
	for key := range raw {
		switch v := raw[key].(type) {
		case string:
			raw[key], err = parseSpecial(v)
			if err != nil {
				return err
			}
		case []any:
			for i := range v {
				if vv, ok := v[i].(string); ok {
					v[i], err = parseSpecial(vv)
					if err != nil {
						return err
					}
				}
			}
		}
	}
	// step 3: assign back
	*a = raw
	return nil
}

func parseSpecial(in string) (any, error) {
	inb := []byte(`"` + in + `"`)
	if strings.HasPrefix(in, "DeviceId(") {
		out := new(DeviceId)
		err := out.UnmarshalJSON(inb)
		return *out, err
	} else if strings.HasPrefix(in, "DeviceClass(") {
		out := new(DeviceClass)
		err := out.UnmarshalJSON(inb)
		return *out, err
	} else if strings.HasPrefix(in, "ChannelType(") {
		out := new(ChannelType)
		err := out.UnmarshalJSON(inb)
		return *out, err
	} else {
		return in, nil
	}
}

type TemplatePayload struct {
	WithPrev bool
	Message  Message
	Queued   []Message
}
