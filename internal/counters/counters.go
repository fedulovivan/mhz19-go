package counters

import (
	"fmt"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type Key = string
type CntrsData = map[Key]utils.Seq
type MinMaxData = map[Key]([2]time.Duration)

const (
	MESSAGES_HANDLED Key = "messagesHandled"
	API_REQUESTS     Key = "apiRequests"
	ERRORS_ALL       Key = "errorsAll"
	QUERIES          Key = "queries"
)

var cntrs = CntrsData{
	MESSAGES_HANDLED: utils.NewSeq(0),
	API_REQUESTS:     utils.NewSeq(0),
	ERRORS_ALL:       utils.NewSeq(0),
	QUERIES:          utils.NewSeq(0),
}
var cntrsMu sync.RWMutex

var minmax = make(MinMaxData)
var minmaxMu sync.RWMutex

func Counters() CntrsData {
	cntrsMu.RLock()
	defer cntrsMu.RUnlock()
	return cntrs
}

func MinMax() MinMaxData {
	minmaxMu.RLock()
	defer minmaxMu.RUnlock()
	return minmax
}

func IncRule(ruleId int) {
	Inc(fmt.Sprintf(
		"rule-%d",
		ruleId,
	))
}

func Time(t time.Duration, key Key) {
	minmaxMu.Lock()
	defer minmaxMu.Unlock()
	// if minmax == nil {
	// 	minmax =
	// }
	c, exist := minmax[key]
	if exist {
		if t < c[0] {
			c[0] = t
		}
		if t > c[1] {
			c[1] = t
		}
		minmax[key] = c
	} else {
		minmax[key] = [2]time.Duration{
			t,
			t,
		}
	}
}

func Inc(key Key) {
	cntrsMu.Lock()
	defer cntrsMu.Unlock()
	c, exist := cntrs[key]
	if !exist {
		c = utils.NewSeq(0)
		cntrs[key] = c
	}
	c.Inc()
}
