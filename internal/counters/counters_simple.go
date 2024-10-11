package counters

import (
	"strconv"
	"sync"
)

type Key string

const (
	MESSAGES_HANDLED Key = "messagesHandled"
	API_REQUESTS     Key = "apiRequests"
	ERRORS_ALL       Key = "errors"
	QUERIES          Key = "queries"
	TRANSACTIONS     Key = "transactions"
)

type Data map[Key]int32

type Container struct {
	sync.Mutex
	data Data
}

func NewContainer(cnt int) *Container {
	return &Container{
		data: make(Data, cnt),
	}
}

var instance = NewContainer(100).Set(API_REQUESTS, 0).Set(ERRORS_ALL, 0)

func (c *Container) Set(key Key, value int32) *Container {
	c.Lock()
	defer c.Unlock()
	c.data[key] = value
	return c
}

func (c *Container) Inc(key Key) {
	c.Lock()
	defer c.Unlock()
	if _, exist := c.data[key]; !exist {
		c.data[key] = 0
	}
	c.data[key]++
}

func (c *Container) Counters() (res Data) {
	c.Lock()
	defer c.Unlock()
	res = make(Data, len(c.data))
	for k, v := range c.data {
		res[k] = v
	}
	return
}

func Counters() Data {
	return instance.Counters()
}

func IncRule(ruleId int) {
	instance.Inc(Key("rule-" + strconv.Itoa(ruleId)))
}

func Inc(key Key) {
	instance.Inc(key)
}
