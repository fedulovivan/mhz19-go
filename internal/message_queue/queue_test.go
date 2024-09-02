package message_queue

import (
	"fmt"
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type QueueSuite struct {
	suite.Suite
}

func (s *QueueSuite) SetupSuite() {
}

func (s *QueueSuite) TeardownSuite() {
}

func (s *QueueSuite) Test10() {
	q := NewQueue(time.Millisecond*100, func(mm []types.Message) {
		fmt.Printf("flushed with %v messages\n", len(mm))
	})
	q.PushMessage(types.Message{})
	q.PushMessage(types.Message{})
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 110)
	q.PushMessage(types.Message{})
	q.PushMessage(types.Message{})
	s.Equal(q.Cnt(), int64(1)) // not yet flushed for a second time
	time.Sleep(time.Millisecond * 110)
	s.Equal(q.Cnt(), int64(2)) // here flushed twice
}

func (s *QueueSuite) Test20() {
	q := NewQueue(time.Millisecond*100, func(mm []types.Message) {
		fmt.Printf("flushed with %v messages\n", len(mm))
	})
	q.PushMessage(types.Message{})
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 50)
	s.Equal(q.Cnt(), int64(0))
	q.PushMessage(types.Message{})
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 50)
	s.Equal(q.Cnt(), int64(0))
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 110)
	s.Equal(q.Cnt(), int64(1))
}

func TestQueue(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}
