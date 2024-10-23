package utils

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type TagSuite struct {
	suite.Suite
}

func (s *TagSuite) Test05() {
	tag1 := NewTag("[main]")
	s.Equal("[main] message one", tag1.F("message one"))
	s.Equal("[main] message two", tag1.F("message two"))
}

func (s *TagSuite) Test10() {
	tag1 := NewTag("[main]")
	tag2 := NewTag("[module]")
	s.Equal("[main] message one", tag1.F("message one"))
	s.Equal("[module] message two", tag2.F("message two"))
}

func (s *TagSuite) Test20() {
	tag1 := NewTag("[main]")
	tag2 := NewTag("[module]")
	tag11 := tag1.With("Rule1")
	tag22 := tag2.With("bar2")
	s.NotEqual(tag1, tag11)
	s.NotEqual(tag2, tag22)
	s.Equal("[main] one", tag1.F("one"))
	s.Equal("[module] two", tag2.F("two"))
	s.Equal("[main] Rule1 one1", tag11.F("one1"))
	s.Equal("[module] bar2 two1", tag22.F("two1"))
}

func (s *TagSuite) Test30() {
	tag1 := NewTag("[main]")
	tag11 := tag1.With("Rule1")
	tag111 := tag11.WithTid("Tid")
	s.NotEqual(tag1, tag11)
	s.NotEqual(tag11, tag111)
	s.NotEqual(tag1, tag111)
	s.Equal("[main] message one", tag1.F("message one"))
	s.Equal("[main] Rule1 rule message", tag11.F("rule message"))
	s.Equal("[main] Rule1 Tid#1 rule and tid message", tag111.F("rule and tid message"))
}

func (s *TagSuite) Test40() {
	s.Equal(
		"foo bar baz message",
		NewTag("foo").With("bar").With("baz").F("message"),
	)
}

func (s *TagSuite) Test50() {
	base := NewTag("[main]")
	tag1 := base.WithTid("Bar")
	tag2 := base.WithTid("Bar")
	tag3 := base.WithTid("Foo")
	s.Equal("[main] Bar#1 msg1", tag1.F("msg1"))
	s.Equal("[main] Bar#2 msg2", tag2.F("msg2"))
	s.Equal("[main] Foo#1 msg3", tag3.F("msg3"))
}

func (s *TagSuite) Test60() {

	s.T().Skip()

	base := NewTag("[main]")
	res := make([]string, 5)
	mu := new(sync.Mutex)

	for i := 0; i < 5; i++ {
		go func(i int) {
			mu.Lock()
			defer mu.Unlock()
			atag := base.With(fmt.Sprintf("action=%d", i))
			res[i] = atag.F("test %d", i)
		}(i)
	}
	time.Sleep(time.Millisecond * 10)
	s.Equal("[main] action=0 test 0", res[0])
	s.Equal("[main] action=1 test 1", res[1])
	s.Equal("[main] action=2 test 2", res[2])
	s.Equal("[main] action=3 test 3", res[3])
	s.Equal("[main] action=4 test 4", res[4])
}

func TestTag(t *testing.T) {
	suite.Run(t, new(TagSuite))
}
