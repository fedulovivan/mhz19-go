package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type RulesSuite struct {
	suite.Suite
}

func (s *RulesSuite) SetupSuite() {
}

func (s *RulesSuite) TeardownSuite() {
}

func (s *RulesSuite) Test10() {
	rule := Rule{
		Condition: Condition{
			Fn: COND_EQUAL,
		},
		Actions: []Action{
			{
				Fn: ACTION_TELEGRAM_BOT_MESSAGE,
			},
		},
	}
	actual, err := json.Marshal(rule)
	s.Nil(err)
	s.Equal(`{"id":0,"condition":{"fn":"Equal"},"actions":[{"fn":"TelegramBotMessage"}],"throttle":null}`, string(actual))
}

func (s *RulesSuite) Test20() {
	cond := Condition{
		Fn: COND_EQUAL,
	}
	actual, err := json.Marshal(cond)
	s.Nil(err)
	s.Equal(`{"fn":"Equal"}`, string(actual))
}

func (s *RulesSuite) Test30() {
	action := Action{
		Fn: ACTION_TELEGRAM_BOT_MESSAGE,
	}
	actual, err := json.Marshal(action)
	s.Nil(err)
	s.Equal(`{"fn":"TelegramBotMessage"}`, string(actual))
}

func (s *RulesSuite) Test40() {
	th := Throttle{}
	actual, err := json.Marshal(th)
	s.Nil(err)
	s.Equal(`null`, string(actual))
}

func (s *RulesSuite) Test41() {
	th := Throttle{Value: time.Minute * 5}
	actual, err := json.Marshal(th)
	s.Nil(err)
	s.Equal(`"5m0s"`, string(actual))
}

func TestRules(t *testing.T) {
	suite.Run(t, new(RulesSuite))
}
