package types

import (
	"fmt"
	"time"
)

// in seconds
// when 0 no "Have not seen" messages will be delivered fot this device
// when NULL a default value DefaultBuriedTimeout/BURIED_TIMEOUT will be used (90m)
// when >0 customised timeout in seconds
type BuriedTimeout struct {
	time.Duration
}

func (d BuriedTimeout) MarshalJSON() ([]byte, error) {
	// if d == nil || d.Duration == 0 {
	// 	return []byte("null"), nil
	// }
	return []byte(fmt.Sprintf(`"%s"`, d)), nil
}
