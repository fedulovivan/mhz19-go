package utils

import "sync"

type Seq interface {
	Inc() int
	Value() int
}

type sequence struct {
	sync.RWMutex
	value int
}

func NewSeq(start int) Seq {
	return &sequence{value: start}
}

func (a *sequence) Value() int {
	a.RLock()
	defer a.RUnlock()
	return a.value
}

func (a *sequence) Inc() (res int) {
	a.Lock()
	defer a.Unlock()
	a.value++
	res = a.value
	return
}
