package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Throttle struct {
	Value time.Duration
}

// https://github.com/golang/go/issues/50480
func (t Throttle) MarshalJSON() (b []byte, err error) {
	if t.Value == 0 {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Value)), nil
}

func (t *Throttle) UnmarshalJSON(b []byte) (err error) {
	var v any
	err = json.Unmarshal(b, &v)
	if err != nil {
		return
	}
	vstring, issting := v.(string)
	if !issting {
		return
	}
	t.Value, err = time.ParseDuration(vstring)
	if err != nil {
		return
	}
	return
}
