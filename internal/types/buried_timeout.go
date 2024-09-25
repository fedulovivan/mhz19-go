package types

import (
	"fmt"
	"time"
)

type BuriedTimeout struct {
	time.Duration
}

func (d BuriedTimeout) MarshalJSON() ([]byte, error) {
	// if d == nil || d.Duration == 0 {
	// 	return []byte("null"), nil
	// }
	return []byte(fmt.Sprintf(`"%s"`, d)), nil
}
