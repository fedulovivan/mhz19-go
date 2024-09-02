package engine

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type ArgReader struct {
	m       types.Message
	args    Args
	mapping Mapping
	errors  []error
}

func (r *ArgReader) Error() error {
	return errors.Join(r.errors...)
}
func (r *ArgReader) Ok() bool {
	return len(r.errors) == 0
}

func (r *ArgReader) Get(field string) any {
	// read argument value
	vRetrieved, err := arg_value(r.m, r.args, field)
	if err != nil {
		r.errors = append(r.errors, err)
	}
	// apply mapping at the end
	if vString, vIsString := vRetrieved.(string); r.mapping != nil && vIsString {
		vMapped, ok := r.mapping[field][vString]
		if ok {
			slog.Debug("ArgReader.Get() %v was mapped to %v", vRetrieved, vMapped)
			return vMapped
		}
	}
	return vRetrieved
}

func NewArgReader(m types.Message, args Args, mapping Mapping) ArgReader {
	return ArgReader{
		m:       m,
		args:    args,
		mapping: mapping,
		errors:  make([]error, 0),
	}
}

func arg_value(m types.Message, args Args, argName string) (any, error) {

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
