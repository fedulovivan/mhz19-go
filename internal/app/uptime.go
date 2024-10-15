package app

import (
	"fmt"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/counters"
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
	ticker := time.NewTicker(time.Second) // update metric each second
	go func() {
		startTime = time.Now()
		for range ticker.C {
			counters.Uptime.Set(time.Since(startTime).Seconds())
		}
	}()
}

func GetUptime() Uptime {
	return Uptime{time.Since(startTime)}
}
