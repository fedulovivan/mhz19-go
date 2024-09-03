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
	m        types.Message
	args     types.Args
	mapping  types.Mapping
	errors   []error
	tpayload *TemplatePayload
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
	value, err := Stage1(r.m, r.args, field)
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
			ttext, err := execTemplate(svalue, r.tpayload)
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
	m types.Message,
	args types.Args,
	mapping types.Mapping,
	tpayload *TemplatePayload,
) ArgReader {
	return ArgReader{
		m:        m,
		args:     args,
		mapping:  mapping,
		errors:   make([]error, 0),
		tpayload: tpayload,
	}
}

func Stage1(m types.Message, args types.Args, argName string) (any, error) {

	// check such arg exist
	v, ok := args[argName]
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
		vg, err := m.Get(field)
		if err != nil {
			return nil, err
		}
		return vg, nil
	}

	if vstring == "$deviceId" {
		return m.DeviceId, nil
	}

	if vstring == "$deviceClass" {
		return m.DeviceClass, nil
	}

	return vstring, nil
}

func isTemplate(in string) bool {
	return strings.Contains(in, "{{") && strings.Contains(in, "}}")
}

func execTemplate(in string, payload *TemplatePayload) (string, error) {
	tmpl, err := template.New("test").Parse(in)
	if err != nil {
		return in, err
	}
	out := bytes.NewBufferString("")
	err = tmpl.Execute(out, payload)
	if err != nil {
		return in, err
	} else {
		return out.String(), nil
	}
}
