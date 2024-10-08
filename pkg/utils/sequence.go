package utils

import (
	"fmt"
	"sync"
)

type Seq interface {
	Inc() int32
	Value() int32
}

type sequence struct {
	sync.RWMutex
	value int32
}

func NewSeq(start int32) Seq {
	return &sequence{value: start}
}

func (a *sequence) MarshalJSON() (b []byte, err error) {
	return []byte(fmt.Sprintf(`%d`, a.Value())), nil
}

func (a *sequence) Value() (res int32) {
	a.RLock()
	res = a.value
	defer a.RUnlock()
	return
}

func (a *sequence) Inc() (res int32) {
	a.Lock()
	defer a.Unlock()
	a.value++
	res = a.value
	return
}
