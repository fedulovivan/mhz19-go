package types

import (
	"encoding/json"
	"fmt"
	"time"
)

type Throttle struct {
	time.Duration
}

// https://github.com/golang/go/issues/50480
func (t Throttle) MarshalJSON() (b []byte, err error) {
	if t.Duration == 0 {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, t.Duration)), nil
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
	t.Duration, err = time.ParseDuration(vstring)
	if err != nil {
		return
	}
	return
}
