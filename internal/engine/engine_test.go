package engine

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
	e   types.Engine
	tag logger.Tag
}

func (s *EngineSuite) SetupSuite() {
	s.e = NewEngine()
	s.e.SetLdmService(ldm.NewService(ldm.RepoSingleton()))
}

var dummy_mtcb types.MessageTupleFn = func(types.DeviceId) (res types.MessageTuple) {
	return
}

func (s *EngineSuite) Test10() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{}, s.tag)
	s.False(actual)
}

func (s *EngineSuite) Test11() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{
		Or: true,
		Nested: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": false}},
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1.11, "Right": 1.11}},
		},
	}, s.tag)
	s.True(actual)
}

func (s *EngineSuite) Test12() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{
		Nested: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": false}},
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1.11, "Right": 1.11}},
		},
	}, s.tag)
	s.False(actual)
}

func (s *EngineSuite) Test13() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{
		Nested: []types.Condition{
			{
				Fn:   types.COND_IS_NIL,
				Args: types.Args{"Value": nil},
			},
			{
				Nested: []types.Condition{
					{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": true}},
					{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
				},
			},
		},
	}, s.tag)
	s.True(actual)
}

func (s *EngineSuite) Test20() {
	defer func() { _ = recover() }()
	s.False(s.e.InvokeConditionFunc(types.MessageTuple{}, 0, false, nil, s.tag))
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test30() {
	actual := s.e.MatchesListSome(dummy_mtcb, []types.Condition{}, s.tag)
	s.False(actual)
}

func (s *EngineSuite) Test31() {
	actual := s.e.MatchesListSome(dummy_mtcb, []types.Condition{
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": "foo", "Right": "bar"}},
	}, s.tag)
	s.True(actual)
}

func (s *EngineSuite) Test40() {
	actual := s.e.MatchesListEvery(dummy_mtcb, []types.Condition{}, s.tag)
	s.False(actual)
}

func (s *EngineSuite) Test41() {
	actual := s.e.MatchesListEvery(dummy_mtcb, []types.Condition{
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": "foo", "Right": "foo"}},
	}, s.tag)
	s.True(actual)
}

func (s *EngineSuite) Test60() {
	s.e.ExecuteActions([]types.Message{}, types.Rule{}, s.tag)
}

func (s *EngineSuite) Test70() {
	// defer func() { _ = recover() }()
	s.e.HandleMessage(types.Message{}, []types.Rule{})
	s.e.HandleMessage(types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}, []types.Rule{})
	// s.Fail("expected to panic")
}

func (s *EngineSuite) Test72() {
	s.e.HandleMessage(types.Message{}, []types.Rule{
		{
			Condition: types.Condition{
				Fn:   types.COND_EQUAL,
				Args: types.Args{"Left": true, "Right": true},
			},
		},
	})
}

func (s *EngineSuite) Test140() {
	s.e.Start()
}

func (s *EngineSuite) Test141() {
	s.e.Stop()
}

func (s *EngineSuite) Test160() {
	input := []byte(`{"Foo":1}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
}

func (s *EngineSuite) Test161() {
	input := []byte(`{"Foo":"bar"}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
	s.Equal(args["Foo"], "bar")
}

func (s *EngineSuite) Test162() {
	input := []byte(`{"Lorem":"DeviceId(bar-111)"}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Lorem")
	s.Equal(args["Lorem"], types.DeviceId("bar-111"))
	s.IsType(args["Lorem"], types.DeviceId(""))
	fmt.Println(args)
}

func (s *EngineSuite) Test163() {
	input := []byte(`{"ClassesList":["DeviceClass(1)","DeviceClass(2)"]}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "ClassesList")
	s.Len(args["ClassesList"], 2)
	s.Equal("map[ClassesList:[zigbee-device (id=1) device-pinger (id=2)]]", fmt.Sprintf("%v", args))
}

func (s *EngineSuite) Test164() {
	input := []byte(`{foo}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	fmt.Println(args)
	s.NotNil(err)
}

func (s *EngineSuite) Test165() {
	input := []byte(`{"foo":"DeviceId"}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Equal("DeviceId", args["foo"])
}

func (s *EngineSuite) Test170() {
	args := types.Args{"Foo1": types.DEVICE_CLASS_BOT}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo1":"telegram-bot"}`, string(argsjson))
}

func (s *EngineSuite) Test171() {
	args := types.Args{"Foo2": types.DeviceId("some-111")}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo2":"some-111"}`, string(argsjson))
}

func TestEngine(t *testing.T) {
	suite.Run(t, &EngineSuite{
		tag: logger.NewTag("EngineSuite"),
	})
}
