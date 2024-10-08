package counters

import (
	"fmt"
	"sync"
)

type Key = string

const (
	MESSAGES_HANDLED Key = "messagesHandled"
	API_REQUESTS     Key = "apiRequests"
	ERRORS_ALL       Key = "errors"
	QUERIES          Key = "queries"
	TRANSACTIONS     Key = "transactions"
)

type Container struct {
	sync.Mutex
	data map[Key]int32
}

var instance = &Container{
	data: map[Key]int32{
		API_REQUESTS: 0,
		ERRORS_ALL:   0,
	},
}

func (c *Container) Inc(key Key) {
	c.Lock()
	defer c.Unlock()
	if _, exist := c.data[key]; !exist {
		c.data[key] = 0
	}
	c.data[key]++
}

func (c *Container) Counters() (res map[Key]int32) {
	c.Lock()
	defer c.Unlock()
	res = make(map[Key]int32, len(c.data))
	for k, v := range c.data {
		res[k] = v
	}
	return
}

func Counters() map[Key]int32 {
	return instance.Counters()
}

func IncRule(ruleId int) {
	instance.Inc(fmt.Sprintf(
		"rule-%d",
		ruleId,
	))
}

func Inc(key Key) {
	instance.Inc(key)
}
