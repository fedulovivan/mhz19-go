package message_queue

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var tag = utils.NewTag(logger.MAIN)

type container struct {
	sync.RWMutex
	qlist map[Key]*queue
}

type Key struct {
	DeviceClass types.DeviceClass
	DeviceId    types.DeviceId
	RuleId      int
}

func (k Key) String() string {
	return fmt.Sprintf("%v-%v-Rule%v", k.DeviceClass, k.DeviceId, k.RuleId)
}

func NewContainer() *container {
	return &container{
		qlist: make(map[Key]*queue),
	}
}

func NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId, ruleId int) (key Key) {
	return Key{
		DeviceClass: deviceClass,
		DeviceId:    deviceId,
		RuleId:      ruleId,
	}
}

func (c *container) HasQueue(key Key) bool {
	c.RLock()
	defer c.RUnlock()
	_, flag := c.qlist[key]
	return flag
}

func (c *container) CreateQueue(key Key, throttle time.Duration, flush OnFlushed) (q *queue) {
	c.Lock()
	defer c.Unlock()
	q = NewQueue(throttle, flush)
	c.qlist[key] = q
	slog.Debug(tag.F(
		"New queue created for key='%v', total instances %v",
		key,
		len(c.qlist),
	))
	return
}

func (c *container) GetQueue(key Key) *queue {
	c.RLock()
	defer c.RUnlock()
	return c.qlist[key]
}
