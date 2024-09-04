package arg_reader

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"strings"
	"text/template"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type ArgReader struct {
	args     types.Args
	message  *types.Message
	mapping  types.Mapping
	errors   []error
	tpayload *TemplatePayload
	engine   types.Engine
}

type TemplatePayload struct {
	Messages []types.Message
}

func (r *ArgReader) Error() error {
	return errors.Join(r.errors...)
}
func (r *ArgReader) Ok() bool {
	return len(r.errors) == 0
}

func (r *ArgReader) Get(field string) any {

	// stage 1: read args map, check field exist, parse special directives ($message, $deviceId, $deviceClass)
	value, err := r.Stage1(field)
	if err != nil {
		r.errors = append(r.errors, err)
	}

	// if value is string go further
	if svalue, isstring := value.(string); isstring {

		// stage 2: apply simple mapping
		if r.mapping != nil {
			vmapped, ok := r.mapping[field][svalue]
			if ok {
				slog.Debug("ArgReader.Get() %v was mapped to %v", value, vmapped)
				value = vmapped
			}
		}

		// stage 3: process arg value as template
		if isTemplate(svalue) && r.tpayload != nil {
			ttext, err := r.ExecTemplate(svalue, field)
			if err != nil {
				r.errors = append(r.errors, err)
			} else {
				value = ttext
			}

		}
	}

	return value
}

func NewArgReader(
	message *types.Message,
	args types.Args,
	mapping types.Mapping,
	tpayload *TemplatePayload,
	engine types.Engine,
) ArgReader {
	return ArgReader{
		message:  message,
		args:     args,
		mapping:  mapping,
		errors:   make([]error, 0),
		tpayload: tpayload,
		engine:   engine,
	}
}

func (r *ArgReader) Stage1(argName string) (any, error) {

	// check such arg exist
	v, ok := r.args[argName]
	if !ok {
		return nil, fmt.Errorf("no such argument %v", argName)
	}

	// return non-strings as is
	vstring, ok := v.(string)
	if !ok {
		return v, nil
	}

	// parse directive
	if strings.HasPrefix(vstring, "$message.") {
		_, field, _ := strings.Cut(vstring, ".")
		vg, err := r.message.Get(field)
		if err != nil {
			return nil, err
		}
		return vg, nil
	}

	if vstring == "$deviceId" {
		return r.message.DeviceId, nil
	}

	if vstring == "$deviceClass" {
		return r.message.DeviceClass, nil
	}

	return vstring, nil
}

func isTemplate(in string) bool {
	return strings.Contains(in, "{{") && strings.Contains(in, "}}")
}

func (r *ArgReader) ExecTemplate(in string, field string) (string, error) {
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
