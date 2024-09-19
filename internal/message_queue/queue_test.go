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
	s.Equal(q.Flushes(), int64(1)) // not yet flushed for a second time
	time.Sleep(time.Millisecond * 110)
	s.Equal(q.Flushes(), int64(2)) // here flushed twice
}

func (s *QueueSuite) Test30() {
	q := NewQueue(time.Millisecond*100, func(mm []types.Message) {
		fmt.Printf("flushed with %v messages\n", len(mm))
	})
	q.PushMessage(types.Message{})
	q.PushMessage(types.Message{})
	s.Equal(q.Flushes(), int64(0))
	time.Sleep(time.Millisecond * 110)
	s.Equal(q.Flushes(), int64(1))
}

func (s *QueueSuite) Test40() {
	flushStats := make([]int, 0, 2)
	q := NewQueue(time.Millisecond*100, func(mm []types.Message) {
		flushStats = append(flushStats, len(mm))
		fmt.Printf("flushed with %v messages\n", len(mm))
	})
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 55)
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 55)
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 55)
	q.PushMessage(types.Message{})
	time.Sleep(time.Millisecond * 55)
	s.Equal(q.Flushes(), int64(2))
	s.Equal([]int{2, 2}, flushStats)
}

func TestQueue(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}
