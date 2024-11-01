package engine

import (
	"io"
	"log"
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
	e types.Engine
}

func (s *EngineSuite) SetupSuite() {
	s.e = NewEngine()
	ldmService := ldm.NewService(ldm.RepoSingleton())
	s.e.SetLdmService(ldmService)
	go func() {
		for range ldmService.OnSet() {
			// noop
			// just to allow send to unbuffered "onset" chan in internal/entities/ldm/repository.go
		}
	}()
}

var dummy_mtcb types.GetCompoundForOtherDeviceId = func(types.DeviceId) (res types.MessageCompound) {
	return
}

func (s *EngineSuite) Test10() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{}, BaseTag)
	s.False(actual)
}

func (s *EngineSuite) Test11() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{
		Or: true,
		Nested: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": false}},
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1.11, "Right": 1.11}},
		},
	}, BaseTag)
	s.True(actual)
}

func (s *EngineSuite) Test12() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{
		Nested: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": false}},
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1.11, "Right": 1.11}},
		},
	}, BaseTag)
	s.False(actual)
}

func (s *EngineSuite) Test13() {
	actual := s.e.MatchesCondition(dummy_mtcb, types.Condition{
		Nested: []types.Condition{
			{
				Fn:   types.COND_NIL,
				Args: types.Args{"Value": nil},
			},
			{
				Nested: []types.Condition{
					{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": true}},
					{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
				},
			},
		},
	}, BaseTag)
	s.True(actual)
}

func (s *EngineSuite) Test20() {
	c := types.Condition{
		Fn: 66,
	}
	s.PanicsWithValue("Condition function 66 not yet implemented", func() {
		s.False(s.e.InvokeConditionFunc(types.MessageCompound{}, c, BaseTag))
	})
}

func (s *EngineSuite) Test30() {
	actual := s.e.MatchesListSome(dummy_mtcb, []types.Condition{}, BaseTag)
	s.False(actual)
}

func (s *EngineSuite) Test31() {
	actual := s.e.MatchesListSome(dummy_mtcb, []types.Condition{
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": "foo", "Right": "bar"}},
	}, BaseTag)
	s.True(actual)
}

func (s *EngineSuite) Test40() {
	actual := s.e.MatchesListEvery(dummy_mtcb, []types.Condition{}, BaseTag)
	s.False(actual)
}

func (s *EngineSuite) Test41() {
	actual := s.e.MatchesListEvery(dummy_mtcb, []types.Condition{
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": "foo", "Right": "foo"}},
	}, BaseTag)
	s.True(actual)
}

func (s *EngineSuite) Test60() {
	s.e.ExecuteActions(types.MessageCompound{}, types.Rule{}, BaseTag)
}

func (s *EngineSuite) Test61() {
	mc := types.MessageCompound{}
	rule := types.Rule{
		Actions: []types.Action{
			{
				Fn: types.ACTION_YEELIGHT_DEVICE_SET_POWER,
			},
			{
				Fn: types.ACTION_PLAY_ALERT,
			},
		},
	}
	tag := utils.NewTag("[Test61]")
	s.e.ExecuteActions(mc, rule, tag)
}

func (s *EngineSuite) Test70() {
	s.e.HandleMessage(types.Message{
		Id:        types.MessageIdSeq.Add(1),
		Timestamp: time.Now(),
	}, []types.Rule{})
}

func (s *EngineSuite) Test71() {
	s.e.HandleMessage(types.Message{
		Id:          types.MessageIdSeq.Add(1),
		Timestamp:   time.Now(),
		DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE,
	}, []types.Rule{})
}

func (s *EngineSuite) Test72() {
	s.e.HandleMessage(types.Message{
		Id:        types.MessageIdSeq.Add(1),
		Timestamp: time.Now(),
	}, []types.Rule{
		{
			Condition: types.Condition{
				Fn:   types.COND_EQUAL,
				Args: types.Args{"Left": true, "Right": true},
			},
		},
	})
}

func Benchmark10(b *testing.B) {

	log.SetOutput(io.Discard)

	e := NewEngine()
	ldmService := ldm.NewService(ldm.RepoSingleton())
	e.SetLdmService(ldmService)
	// noop
	// just to allow send to unbuffered "onset" chan in internal/entities/ldm/repository.go
	go func() {
		for range ldmService.OnSet() {
		}
	}()

	for k := 0; k < b.N; k++ {
		e.HandleMessage(types.Message{
			Id:        types.MessageIdSeq.Add(1),
			Timestamp: time.Now(),
		}, []types.Rule{
			// {
			// 	Condition: types.Condition{
			// 		Fn:   types.COND_EQUAL,
			// 		Args: types.Args{"Left": true, "Right": true},
			// 	},
			// },
			// {
			// 	Condition: types.Condition{
			// 		Or: true,
			// 		Nested: []types.Condition{
			// 			{Fn: types.COND_CHANGED},
			// 			{Fn: types.COND_EQUAL},
			// 			{Fn: types.COND_IN_LIST},
			// 			{Fn: types.COND_NIL},
			// 			{Fn: types.COND_ZIGBEE_DEVICE},
			// 			{Fn: types.COND_DEVICE_CLASS},
			// 			{Fn: types.COND_СHANNEL},
			// 			{Fn: types.COND_FROM_END_DEVICE},
			// 			{Fn: types.COND_TRUE},
			// 			{Fn: types.COND_FALSE},
			// 			{Fn: types.COND_DEVICE_ID},
			// 		},
			// 	},
			// },
		})
	}
}

func (s *EngineSuite) Test140() {
	s.e.Start()
}

func (s *EngineSuite) Test141() {
	s.e.Stop()
}

func TestEngine(t *testing.T) {
	suite.Run(t, &EngineSuite{})
}
