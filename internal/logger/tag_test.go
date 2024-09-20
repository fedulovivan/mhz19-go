package logger

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TagSuite struct {
	suite.Suite
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

func TestTag(t *testing.T) {
	suite.Run(t, new(TagSuite))
}
