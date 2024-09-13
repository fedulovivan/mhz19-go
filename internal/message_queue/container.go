package message_queue

import (
	"fmt"
	"log/slog"
	"strconv"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var tag = logger.NewTag(logger.MAIN)

type container struct {
	mu    sync.RWMutex
	qlist map[Key]Queue
}

type Key string

func NewContainer() *container {
	return &container{
		qlist: make(map[Key]Queue),
	}
}

func (c *container) MakeKey(deviceClass types.DeviceClass, deviceId types.DeviceId, ruleId int) (key Key) {
	return Key(
		types.DEVICE_CLASS_NAMES[deviceClass] + "-" + string(deviceId) + "-rule" + strconv.Itoa(int(ruleId)),
	)
}

func (c *container) HasQueue(key Key) (flag bool) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	_, flag = c.qlist[key]
	return
}

func (c *container) CreateQueue(key Key, throttle time.Duration, flush FlushFn) (q Queue) {
	c.mu.Lock()
	defer c.mu.Unlock()
	q = NewQueue(throttle, flush)
	c.qlist[key] = q
	slog.Debug(tag.F(fmt.Sprintf(
		"New Queue created for key='%v', total instances %v",
		key,
		len(c.qlist),
	)))
	return
}

func (c *container) GetQueue(key Key) (qq Queue) {
	c.mu.RLock()
	defer c.mu.RUnlock()
	return c.qlist[key]
}
