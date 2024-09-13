package logger

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type TagSuite struct {
	suite.Suite
}

func (s *TagSuite) Test1() {
	tag1 := NewTag("[main]")
	tag2 := NewTag("[module]")
	s.Equal("[main] message one", tag1.F("message one"))
	s.Equal("[module] message two", tag2.F("message two"))
}

func (s *TagSuite) Test2() {
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

func (s *TagSuite) Test3() {
	tag1 := NewTag("[main]")
	tag11 := tag1.With("Rule1")
	tag111 := tag11.WithTid()
	s.NotEqual(tag1, tag11)
	s.NotEqual(tag11, tag111)
	s.NotEqual(tag1, tag111)
	s.Equal("[main] message one", tag1.F("message one"))
	s.Equal("[main] Rule1 rule message", tag11.F("rule message"))
	s.Equal("[main] Rule1 Tid#1 rule and tid message", tag111.F("rule and tid message"))
}

func (s *TagSuite) Test4() {
	s.Equal(
		"foo bar baz message",
		NewTag("foo").With("bar").With("baz").F("message"),
	)
}

func TestTag(t *testing.T) {
	suite.Run(t, new(TagSuite))
}
