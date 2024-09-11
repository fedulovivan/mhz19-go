package message_queue

import (
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type FlushFn = func(mm []types.Message)

type Queue interface {
	// store next message
	PushMessage(m types.Message)
	// for uts
	Cnt() int64
}

var _ Queue = (*queue)(nil)

type queue struct {
	cnt      int64
	mu       sync.RWMutex
	mm       []types.Message
	throttle time.Duration
	flushCb  FlushFn
	timer    *time.Timer
}

func (q *queue) onFlushed() {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.flushCb(q.mm)
	q.cnt++
	q.timer = nil
	q.mm = make([]types.Message, 0)
}

func (q *queue) Cnt() int64 {
	q.mu.RLock()
	defer q.mu.RUnlock()
	return q.cnt
}

func (q *queue) PushMessage(m types.Message) {
	q.mu.Lock()
	defer q.mu.Unlock()
	q.mm = append(q.mm, m)
	if q.timer == nil {
		q.timer = time.AfterFunc(q.throttle, q.onFlushed)
	} else {
		q.timer.Reset(q.throttle)
	}
}

func NewQueue(throttle time.Duration, flush FlushFn) Queue {
	return &queue{
		throttle: throttle,
		flushCb:  flush,
	}
}
