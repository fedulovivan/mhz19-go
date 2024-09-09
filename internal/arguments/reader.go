package arguments

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"text/template"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type reader struct {
	args     types.Args
	message  *types.Message
	errors   []error
	mapping  types.Mapping
	tpayload *types.TemplatePayload
	engine   types.EngineAsSupplier
}

func NewReader(
	message *types.Message,
	args types.Args,
	mapping types.Mapping,
	tpayload *types.TemplatePayload,
	engine types.EngineAsSupplier,
) reader {
	return reader{
		message:  message,
		args:     args,
		mapping:  mapping,
		tpayload: tpayload,
		engine:   engine,
		errors:   make([]error, 0),
	}
}

func (r *reader) Error() error {
	return errors.Join(r.errors...)
}

func (r *reader) push_err(err error) {
	r.errors = append(r.errors, err)
}

func (r *reader) Get(field string) any {

	// stage 1: check if requested arg exist in map
	in, exist := r.args[field]
	if !exist {
		r.push_err(fmt.Errorf("no such argument %v", field))
		return nil
	}

	var out any = in

	if sIn, isString := in.(string); isString {
		if isTemplate(sIn) {
			// stage 2: process string value as template
			processed, err := r.ExecTemplate(sIn, field)
			if err == nil {
				out = processed
			} else {
				r.push_err(err)
			}
		} else if r.message.IsSpecial(sIn) {
			// stage 2: process string argument as special Message's getter directive
			processed, err := r.message.ExecDirective(sIn)
			if err == nil {
				out = processed
			} else {
				r.push_err(err)
			}
		}
	}

	// stage 3: as a final step apply mapping
	if r.mapping != nil {
		if fieldMap, ok := r.mapping[field]; ok {
			outAsKey := fmt.Sprintf("%v", out)
			if mapped, ok := fieldMap[outAsKey]; ok {
				slog.Debug(
					"ArgReader.Get() in=%v (out=%v, outAsKey=%v) was mapped to %v",
					in, out, outAsKey, mapped,
				)
				out = mapped
			}
		}
	}

	return out
}

func (r *reader) ExecTemplate(in string, field string) (string, error) {
	tmpl, err := template.New(field).Funcs(
		template.FuncMap{
			"deviceName": func(deviceId any) (string, error) {
				if typedDeviceId, ok := deviceId.(types.DeviceId); ok {
					device, err := r.engine.DevicesService().GetOne(typedDeviceId)
					if err != nil {
						return string(typedDeviceId), err
					}
					if device.Name == "" {
						return fmt.Sprintf("<unknonwn device originated from %v> %v", device.Origin, deviceId), nil
					}
					return device.Name, nil
				} else {
					return fmt.Sprintf("%v", deviceId), fmt.Errorf(
						"deviceName accepts only types.DeviceId as an argument",
					)
				}
			},
			"pingerStatusName": func(statusId any) string {
				svalue := fmt.Sprintf("%v", statusId)
				if svalue == "0" {
					return "OFFLINE"
				} else if svalue == "1" {
					return "ONLINE"
				} else {
					return "UNKNOWN"
				}
			},
		},
	).Parse(in)
	if err != nil {
		return in, err
	}
	out := bytes.NewBufferString("")
	slog.Debug(fmt.Sprintf("executing template '%v' with data %+v", in, r.tpayload))
	err = tmpl.Execute(out, r.tpayload)
	if err != nil {
		return in, err
	} else {
		return out.String(), nil
	}
}

// result = stringValue
// if err == nil {
// 	result = processed
// } else {
// 	r.push_err(err)
// 	result = stringValue
// }

// func (r *reader) ExecMessageFieldGetter(in string) (any, error) {
// 	if r.message == nil {
// 		return in, nil
// 	}
// 	_, field, _ := strings.Cut(in, ".")
// 	return r.message.Get(field)
// }

// TODO move to internal/types/message.go::Get()
// func (r *reader) ExecSimpleDirective(in string) any {
// 	if r.message == nil {
// 		return in
// 	}
// 	if in == "$deviceId" {
// 		return r.message.DeviceId
// 	} else if in == "$deviceClass" {
// 		return r.message.DeviceClass
// 	} else if in == "$channelType" {
// 		return r.message.ChannelType
// 	} else {
// 		panic(fmt.Sprintf("ExecDirective() unknown directive %v", in))
// 	}
// }

// process arg value as template
// }
// return vstring, nil
// stage 1: read args map, check field exist, parse special directives ($message, $deviceId, $deviceClass)
// value, err := r.stage_one(field)
// if err != nil {
// 	r.errors = append(r.errors, err)
// }
// if value is string go further
// if svalue, isstring := value.(string); isstring {
// func (r *reader) stage_one(argName string) (any, error) {
// }
// if isSimpleDirective(stringValue) {
// 	// process special directive for reading message field
// 	result = r.ExecSimpleDirective(stringValue)
// } else if strings.HasPrefix(stringValue, "$message.") {
// 	// process special directive for reading message field
// 	processed, err := r.ExecMessageFieldGetter(stringValue)
// 	if err == nil {
// 		result = processed
// 	} else {
// 		r.push_err(err)
// 		result = stringValue
// 	}
// } else
