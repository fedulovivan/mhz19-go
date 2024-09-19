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
	s.Equal(c.GetQueue(key).Cnt(), int64(0))

	key2 := NewKey(types.DEVICE_CLASS_ZIGBEE_DEVICE, types.DeviceId("foo2"), 111)
	s.False(c.HasQueue(key2))
	s.Nil(c.GetQueue(key2))
}

func TestContainer(t *testing.T) {
	suite.Run(t, new(ContainerSuite))
}
