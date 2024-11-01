package counters

import (
	"fmt"
	"sync"
	"sync/atomic"
	"time"
)

type Metric struct {
	time.Duration
}

func (m *Metric) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`"%s"`, m.Duration)), nil
}

type Counter struct {
	atomic.Int32
}

func (m *Counter) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`%d`, m.Int32.Load())), nil
}

type TimingsRecord struct {
	Cnt Counter `json:"total"`
	Min Metric  `json:"min"`
	Max Metric  `json:"max"`
	Avg Metric  `json:"avg"`
}

type TimingsData map[Key]*TimingsRecord

var timings = make(TimingsData, 0)
var timingsMu sync.RWMutex

func TimingsCopy() (res TimingsData) {
	timingsMu.RLock()
	defer timingsMu.RUnlock()
	res = make(TimingsData, len(timings))
	for k, v := range timings {
		res[k] = &TimingsRecord{
			Min: Metric{Duration: v.Min.Duration},
			Max: Metric{Duration: v.Max.Duration},
			Avg: Metric{Duration: v.Avg.Duration},
			Cnt: Counter{},
		}
		res[k].Cnt.Int32.Add(v.Cnt.Load())
	}
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
			Cnt: Counter{},
		}
		timings[key].Cnt.Add(1)
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
	// TODO: could be simplified: replace with sum and count (as in prometheus)
	avg := (int64(record.Avg.Duration)*int64(record.Cnt.Load()) + int64(d)) / int64(record.Cnt.Add(1))
	record.Avg.Duration = time.Duration(avg)

}
