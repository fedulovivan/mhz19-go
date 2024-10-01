package counters

import (
	"fmt"
	"sync"

	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var data = map[Key]utils.Seq{
	MESSAGES_RECEIVED: utils.NewSeq(0),
	API_REQUESTS:      utils.NewSeq(0),
	ERRORS:            utils.NewSeq(0),
}
var lock sync.Mutex

type Key = string

const (
	MESSAGES_RECEIVED Key = "messagesReceived"
	API_REQUESTS      Key = "apiRequests"
	ERRORS            Key = "errors"
)

func Data() map[Key]utils.Seq {
	return data
}

func IncRule(ruleId int) {
	Inc(fmt.Sprintf(
		"rule-%d",
		ruleId,
	))
}

func Inc(key Key) {
	lock.Lock()
	defer lock.Unlock()
	c, exist := data[key]
	if !exist {
		c = utils.NewSeq(0)
		data[key] = c
	}
	c.Inc()
}
