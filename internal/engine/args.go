package engine

import (
	"errors"
	"fmt"
	"log/slog"
	"strings"
)

type Args map[string]any

type ArgReader struct {
	fn     CondFn
	m      Message
	args   Args
	errors []error
}

func (r *ArgReader) Ok() bool {
	valid := len(r.errors) == 0
	if !valid {
		slog.Error(string(r.fn) + ": " + errors.Join(r.errors...).Error())
	}
	return valid
}

func (r *ArgReader) Get(field string) any {
	v, err := arg_value(r.m, r.args, field)
	if err != nil {
		r.errors = append(r.errors, err)
	}
	return v
}

// func Get[T any](c ArgReader, field string) T {
// 	v, err := arg_value(c.m, c.args, field)
// 	if err != nil {
// 		c.errors = append(c.errors, err)
// 	}
// 	vt, ok := v.(T)
// 	if !ok {
// 		c.errors = append(c.errors, fmt.Errorf("cannot cast"))
// 	}
// 	return vt
// }

func NewArgReader(fn CondFn, m Message, args Args) ArgReader {
	return ArgReader{
		fn,
		m,
		args,
		make([]error, 0),
	}
}

func arg_value(m Message, args Args, argName string /* , ee *[]error */) (any, error) {

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

// func EqualFn
// func Coerce[T any](mt MessageTuple, args Args, name string, errors *[]error) T {
// 	v, err := value(mt[0], args, name)
// 	if err != nil {
// 		*errors = append(*errors, err)
// 	}
// 	vt, _ := v.(T)
// 	return vt
// 	// return v
// 	// if ok {
// 	// 	*errors = append(*errors, fmt.Errorf("unexpected type %T for %v", v, v))
// 	// }
// 	// return vt
// }
// func Coerce[T any](mt MessageTuple, args Args, name string, errors *[]error) T {
// 	v, err := value(mt[0], args, name)
// 	if err != nil {
// 		*errors = append(*errors, err)
// 	}
// 	vt, ok := v.(T)
// 	if ok {
// 		*errors = append(*errors, fmt.Errorf("unexpected type %T for %v", v, v))
// 	}
// 	return vt
// }
// left := args.Get[string](mt, "Left")
// left := args.Get[string](mt, "Left")
// return EqualFnGen[string](
// )
// g, ctx := errgroup.WithContext(context.Background())
// var errors []error
// fmt.Println("errors1", errors)
// fmt.Println("errors2", errors)
// err := g.Wait()
// v := Coerce[any](mt, args, "Left")
// left, e1 := value(mt[0], args, "Left")
// right, e2 := value(mt[0], args, "Right")
// if e1 != nil || e2 != nil {
// 	return false
// }
// return left == right
