package message_queue

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type OnFlushed func(mm []types.Message)

type queue struct {
	sync.RWMutex
	flushesCount atomic.Int64
	waitFlush    atomic.Pointer[chan struct{}]
	mm           []types.Message
	throttle     time.Duration
	onFlushed    OnFlushed
	timer        *time.Timer
}

func (q *queue) flush_internal() {
	q.Lock()
	defer q.Unlock()
	if q.onFlushed != nil {
		q.onFlushed(q.mm)
	}
	q.timer = nil
	q.mm = nil
	q.flushesCount.Add(1)
	close(*q.waitFlush.Load())
}

func (q *queue) Flush() {
	q.Lock()
	defer q.Unlock()
	if q.timer != nil {
		q.timer.Reset(0)
	}
}

func (q *queue) Wait() {
	wf := q.waitFlush.Load()
	if wf == nil {
		return
	}
	<-*wf
}

func (q *queue) Flushes() int64 {
	return q.flushesCount.Load()
}

func (q *queue) PushMessage(m types.Message) {
	q.Lock()
	defer q.Unlock()
	q.mm = append(q.mm, m)
	if q.timer == nil {
		wf := make(chan struct{})
		q.waitFlush.Store(&wf)
		q.timer = time.AfterFunc(q.throttle, q.flush_internal)
	}
}

func NewQueue(throttle time.Duration, onFlushed OnFlushed) *queue {
	return &queue{
		throttle:  throttle,
		onFlushed: onFlushed,
	}
}
