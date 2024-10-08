package counters

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type CountersSuite struct {
	suite.Suite
}

func (s *CountersSuite) SetupSuite() {
}

func (s *CountersSuite) TeardownSuite() {
}

func (s *CountersSuite) Test10() {
	Inc("foo10")
	data := Counters()
	s.Equal(int32(1), data["foo10"])
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
	s.Equal(int32(1000), data["foo30"])
	s.Equal(int32(1000), data["bar30"])
}

func (s *CountersSuite) Test31() {
	wg := sync.WaitGroup{}
	iterations := 1000
	expected := 2000
	key := "foo31"
	wg.Add(3)
	go func() {
		for i := 0; i < iterations; i++ {
			Inc(key)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < iterations; i++ {
			Inc(key)
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < iterations; i++ {
			data := Counters()
			if v, ok := data[key]; ok {
				json.Marshal(v)
			}
		}
		wg.Done()
	}()
	wg.Wait()
	s.Equal(int32(expected), Counters()[key])
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
	s.Equal(int32(1), data["rule-1"])
	s.Equal(int32(1), data["rule-2"])
}

func (s *CountersSuite) Test60() {
	Time(time.Millisecond, "foo-1")
	Time(time.Millisecond*10, "foo-1")
	Time(time.Millisecond*100, "foo-1")
	data := Timings()
	jdata, _ := json.Marshal(data)
	s.Contains(string(jdata), `"foo-1":{"total":3,"min":"1ms","max":"100ms","avg":"37ms"}`)
}

func (s *CountersSuite) Test61() {
	Time(time.Millisecond*2, "foo-2")
	Time(time.Millisecond*3, "foo-2")
	Time(time.Millisecond*1, "foo-2")
	Time(time.Millisecond*5, "foo-2")
	Time(time.Millisecond*4, "foo-2")
	data := Timings()
	jdata, _ := json.Marshal(data)
	s.Contains(string(jdata), `"foo-2":{"total":5,"min":"1ms","max":"5ms","avg":"3ms"}`)
}

func (s *CountersSuite) Test62() {
	for i := 0; i < 1000; i++ {
		Time(time.Millisecond, "foo-3")
	}
	data := Timings()
	jdata, _ := json.Marshal(data["foo-3"])
	s.Equal(`{"total":1000,"min":"1ms","max":"1ms","avg":"1ms"}`, string(jdata))
}

func (s *CountersSuite) Test70() {
	iterations := 1
	Time(time.Millisecond, "foo-4")
	wg := sync.WaitGroup{}
	wg.Add(3)
	go func() {
		for i := 0; i < iterations; i++ {
			Time(time.Millisecond, "foo-4")
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < iterations; i++ {
			Time(time.Millisecond, "foo-4")
		}
		wg.Done()
	}()
	go func() {
		for i := 0; i < iterations; i++ {
			data := Timings()
			json.Marshal(data)
		}
		wg.Done()
	}()
	wg.Wait()
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

func Benchmark30(b *testing.B) {
	for k := 0; k < b.N; k++ {
		Time(time.Microsecond, "foo-4")
	}
}

func TestCounters(t *testing.T) {
	suite.Run(t, new(CountersSuite))
}
