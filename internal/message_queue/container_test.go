package message_queue

import (
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type ContainerSuite struct {
	suite.Suite
}

func (s *ContainerSuite) SetupSuite() {
}

func (s *ContainerSuite) TeardownSuite() {
}

func (s *ContainerSuite) Test10() {
	c := NewContainer()

	key := NewKey(types.DEVICE_CLASS_ZIGBEE_DEVICE, types.DeviceId("foo1"), 111)
	s.False(c.HasQueue(key))
	c.CreateQueue(key, time.Millisecond*100, func(mm []types.Message) {})
	s.True(c.HasQueue(key))
	s.NotNil(c.GetQueue(key))
	c.GetQueue(key).PushMessage(types.Message{})
	s.Equal(c.GetQueue(key).Flushes(), int64(0))

	key2 := NewKey(types.DEVICE_CLASS_ZIGBEE_DEVICE, types.DeviceId("foo2"), 111)
	s.False(c.HasQueue(key2))
	s.Nil(c.GetQueue(key2))
}

func (s *ContainerSuite) Test20() {
	// race condition test
	c := NewContainer()
	key := NewKey(types.DEVICE_CLASS_ZIGBEE_DEVICE, types.DeviceId("foo1"), 111)
	done := make(chan int)
	go func() {
		for i := 0; i < 1000; i++ {
			c.HasQueue(key)
		}
		done <- 1
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			c.CreateQueue(key, 0, nil)
		}
		done <- 1
	}()
	<-done
	<-done
}

// wait for several queues to finish
func (s *ContainerSuite) Test30() {
	c := NewContainer()
	key1 := NewKey(types.DEVICE_CLASS_ZIGBEE_DEVICE, types.DeviceId("foo1"), 111)
	key2 := NewKey(types.DEVICE_CLASS_ZIGBEE_DEVICE, types.DeviceId("foo2"), 111)
	q1 := c.CreateQueue(key1, time.Millisecond*10, nil)
	q2 := c.CreateQueue(key2, time.Millisecond*15, nil)
	q1.PushMessage(types.Message{})
	q2.PushMessage(types.Message{})
	s.Eventually(func() bool {
		c.Wait()
		c.Wait()
		return true
	}, time.Millisecond*20, time.Millisecond*10)
}

// func (s *ContainerSuite) Test40() {
// 	tt := []*time.Timer{
// 		time.NewTimer(time.Second),
// 		time.NewTimer(time.Second + time.Millisecond*100),
// 	}
// 	for _, t := range tt {
// 		<-t.C
// 	}
// }

func TestContainer(t *testing.T) {
	suite.Run(t, new(ContainerSuite))
}
