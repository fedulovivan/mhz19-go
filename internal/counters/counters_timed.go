package counters

import (
	"fmt"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type Metric struct {
	time.Duration
}

func (m *Metric) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, m.Duration)), nil
}

type TimingsRecord struct {
	Cnt utils.Seq `json:"total"`
	Min Metric    `json:"min"`
	Max Metric    `json:"max"`
	Avg Metric    `json:"avg"`
}

type TimingsData = map[Key]*TimingsRecord

var timings = make(TimingsData)
var timingsMu sync.RWMutex

func Timings() (res TimingsData) {
	timingsMu.Lock()
	res = make(TimingsData, len(timings))
	for k, v := range timings {
		res[k] = &TimingsRecord{
			Min: Metric{Duration: v.Min.Duration},
			Max: Metric{Duration: v.Max.Duration},
			Avg: Metric{Duration: v.Avg.Duration},
			Cnt: utils.NewSeq(v.Cnt.Value()),
		}
	}
	defer timingsMu.Unlock()
	return
}

func TimeSince(d time.Time, key Key) {
	Time(time.Since(d), key)
}

func Time(d time.Duration, key Key) {

	timingsMu.Lock()
	defer timingsMu.Unlock()

	record, exist := timings[key]
	if !exist {
		timings[key] = &TimingsRecord{
			Min: Metric{Duration: d},
			Max: Metric{Duration: d},
			Avg: Metric{Duration: d},
			Cnt: utils.NewSeq(1),
		}
		return
	}

	// Min
	if d < record.Min.Duration {
		record.Min.Duration = d
	}

	// Max
	if d > record.Max.Duration {
		record.Max.Duration = d
	}

	// Avg + Total
	// https://math.stackexchange.com/a/4456459
	avg := (int64(record.Avg.Duration)*int64(record.Cnt.Value()) + int64(d)) / int64(record.Cnt.Inc())
	record.Avg.Duration = time.Duration(avg)

}
