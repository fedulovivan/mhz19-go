package message_queue

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type OnFlushed func(mm []types.Message)

type queue struct {
	sync.Mutex
	flushes   atomic.Int64
	mm        []types.Message
	throttle  time.Duration
	onFlushed OnFlushed
	timer     *time.Timer
}

func (q *queue) flush() {
	q.Lock()
	defer q.Unlock()
	q.onFlushed(q.mm)
	q.timer = nil
	q.mm = nil
	q.flushes.Add(1)
}

func (q *queue) Flushes() int64 {
	return q.flushes.Load()
}

func (q *queue) PushMessage(m types.Message) {
	q.Lock()
	defer q.Unlock()
	q.mm = append(q.mm, m)
	if q.timer == nil {
		q.timer = time.AfterFunc(q.throttle, q.flush)
	}
}

func NewQueue(throttle time.Duration, onFlushed OnFlushed) *queue {
	return &queue{
		throttle:  throttle,
		onFlushed: onFlushed,
	}
}
