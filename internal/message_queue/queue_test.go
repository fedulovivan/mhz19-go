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

func (s *QueueSuite) Test50() {
	// for the queue without pushed messages Wait returns immediately
	q := NewQueue(time.Second*5, nil)
	s.Eventually(func() bool {
		q.Wait()
		return true
	}, time.Millisecond*5, time.Millisecond*1)
}

func (s *QueueSuite) Test60() {
	q := NewQueue(time.Millisecond*10, nil)
	q.PushMessage(types.Message{})
	s.Eventually(func() bool {
		q.Wait()
		return true
	}, time.Millisecond*15, time.Millisecond*10)
}

func (s *QueueSuite) Test70() {
	q1 := NewQueue(time.Millisecond*10, nil)
	q2 := NewQueue(time.Millisecond*15, nil)
	q1.PushMessage(types.Message{})
	q2.PushMessage(types.Message{})
	s.Eventually(func() bool {
		q1.Wait()
		q2.Wait()
		return true
	}, time.Millisecond*20, time.Millisecond*5)
}

func (s *QueueSuite) Test80() {
	q := NewQueue(time.Millisecond*100, nil)
	q.PushMessage(types.Message{})
	q.Flush()
	s.Eventually(func() bool {
		q.Wait()
		return true
	}, time.Millisecond*10, time.Millisecond*5)
	s.Equal(q.Flushes(), int64(1))
}

func TestQueue(t *testing.T) {
	suite.Run(t, new(QueueSuite))
}
