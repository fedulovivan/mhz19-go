package arguments

import (
	"bytes"
	"errors"
	"fmt"
	"log/slog"
	"text/template"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// var baseTag = logger.NewTag(logger.ARGS)

type reader struct {
	args     types.Args
	message  *types.Message
	errors   []error
	mapping  types.Mapping
	tpayload *types.TemplatePayload
	engine   types.EngineAsSupplier
	baseTag  logger.Tag
}

func NewReader(
	message *types.Message,
	args types.Args,
	mapping types.Mapping,
	tpayload *types.TemplatePayload,
	engine types.EngineAsSupplier,
	baseTag logger.Tag,
) reader {
	return reader{
		message:  message,
		args:     args,
		mapping:  mapping,
		tpayload: tpayload,
		engine:   engine,
		errors:   make([]error, 0),
		baseTag:  baseTag,
	}
}

func (r *reader) Error() error {
	return errors.Join(r.errors...)
}

func (r *reader) push_err(err error) {
	r.errors = append(r.errors, err)
}

func GetTyped[T any](r *reader, field string) (res T, err error) {
	v := r.Get(field)
	err = r.Error()
	if err != nil {
		return
	}
	// if v == nil {
	// 	return
	// }
	res, ok := v.(T)
	if !ok {
		err = fmt.Errorf("cannot cast %T to %T", v, res)
	}
	return
}

func (r *reader) Has(field string) (has bool) {
	_, has = r.args[field]
	return
}

func (r *reader) Get(field string) any {

	tag := r.baseTag.With(`reader.Get("%s"):`, field)

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
		} else if types.IsSpecialDirective(sIn) {
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
				if app.Config.ArgsDebug {
					slog.Debug(tag.F(
						`in="%v" (out="%v", outAsKey="%v") was mapped to "%v"`,
						in, out, outAsKey, mapped,
					))
				}
				out = mapped
			}
		}
	}

	if app.Config.ArgsDebug {
		slog.Debug(tag.F("in=%v (%T), out=%v (%T)", in, in, out, out))
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
						return string(typedDeviceId), nil
					}
					if device.Name != nil {
						return *device.Name, nil
					}
					if device.Comments != nil {
						return *device.Comments, nil
					}
					return fmt.Sprintf("Device of class %s, with id %v", device.DeviceClassId, deviceId), nil
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
				} else if svalue == "-1" {
					return "UNKNOWN"
				} else {
					return svalue
				}
			},
			"openedClosed": func(contact any) string {
				svalue := fmt.Sprintf("%v", contact)
				if svalue == "1" || svalue == "true" {
					return "closed"
				} else if svalue == "0" || svalue == "false" {
					return "opened"
				} else {
					return svalue
				}
			},
			"leakage": func(leakage bool) string {
				if leakage {
					return "is leaking"
				} else {
					return "is dry"
				}
			},
			"time": func(t time.Time) string {
				return t.Format("15:04:05")
			},
		},
	).Parse(in)
	if err != nil {
		return in, err
	}
	out := bytes.NewBufferString("")
	if app.Config.ArgsDebug {
		slog.Debug(fmt.Sprintf("executing template '%v' with data %+v", in, r.tpayload))
	}
	err = tmpl.Execute(out, r.tpayload)
	if err != nil {
		return in, err
	} else {
		return out.String(), nil
	}
}
