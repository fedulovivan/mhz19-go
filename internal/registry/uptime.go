package registry

import (
	"fmt"
	"time"
)

var startTime time.Time

type Uptime struct {
	time.Duration
}

func (d Uptime) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func RecordStartTime() {
	if !startTime.IsZero() {
		panic("expected to be called only once")
	}
	startTime = time.Now()
}

func GetUptime() Uptime {
	return Uptime{time.Since(startTime)}
}
