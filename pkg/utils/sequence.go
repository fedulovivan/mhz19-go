package utils

import "sync"

type Seq interface {
	Next() (id int)
}

type sequence struct {
	sync.Mutex
	id int
}

func NewSeq() Seq {
	return new(sequence)
}

func (a *sequence) Next() (id int) {
	a.Lock()
	defer a.Unlock()
	a.id++
	id = a.id
	return
}
