package counters

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CountersSuite struct {
	suite.Suite
}

func (s *CountersSuite) Test10() {
	Inc("foo10")
	data := Counters()
	s.Equal(int32(1), data["foo10"].Value())
}

func (s *CountersSuite) Test20() {
	Inc("foo20")
	Inc("bar20")
	Inc("bar20")
	bytes, err := json.Marshal(Counters())
	s.Nil(err)
	s.Contains(string(bytes), `"bar20":2`)
	s.Contains(string(bytes), `"foo20":1`)
}

func (s *CountersSuite) Test30() {
	for i := 0; i < 1000; i++ {
		Inc("foo30")
		Inc("bar30")
	}
	data := Counters()
	s.Equal(int32(1000), data["foo30"].Value())
	s.Equal(int32(1000), data["bar30"].Value())
}

func (s *CountersSuite) Test31() {
	go func() {
		for i := 0; i < 1000; i++ {
			Inc("foo31")
		}
	}()
	go func() {
		for i := 0; i < 1000; i++ {
			Inc("foo31")
		}
	}()
	data := Counters()
	time.Sleep(time.Millisecond * 100)
	s.Equal(int32(2000), data["foo31"].Value())
}

func (s *CountersSuite) Test40() {
	s.NotPanics(func() {
		data := Counters()
		fmt.Println(data)
	})
}

func (s *CountersSuite) Test50() {
	IncRule(1)
	IncRule(2)
	data := Counters()
	s.Equal(int32(1), data["rule-1"].Value())
	s.Equal(int32(1), data["rule-2"].Value())
}

func Benchmark10(b *testing.B) {
	for k := 0; k < b.N; k++ {
		Inc("lorem10")
	}
}

func Benchmark20(b *testing.B) {
	for k := 0; k < b.N; k++ {
		Inc(fmt.Sprintf("%v", k))
	}
}

func TestCounters(t *testing.T) {
	suite.Run(t, new(CountersSuite))
}
