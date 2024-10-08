package utils

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/suite"
)

type SeqSuite struct {
	suite.Suite
}

func (s *SeqSuite) SetupSuite() {
}

func (s *SeqSuite) TeardownSuite() {
}

func (s *SeqSuite) Test10() {
	seq := NewSeq(0)
	seq.Inc()
	seq.Value()
}

func (s *SeqSuite) Test20() {
	wg := sync.WaitGroup{}
	seq := NewSeq(0)
	iterations := 100
	expected := int32(200)
	wg.Add(3)
	go func() {
		for i := 0; i < iterations; i++ {
			seq.Inc()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < iterations; i++ {
			seq.Inc()
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < iterations; i++ {
			seq.Value()
		}
		wg.Done()
	}()
	wg.Wait()
	s.Equal(expected, seq.Value())
}

func TestSeq(t *testing.T) {
	suite.Run(t, new(SeqSuite))
}
