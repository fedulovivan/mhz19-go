package actions

import (
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/mocks"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type ActionsSuite struct {
	suite.Suite
}

func (s *ActionsSuite) SetupSuite() {
}

func (s *ActionsSuite) TeardownSuite() {
}

func (s *ActionsSuite) Test10() {
	engine := mocks.NewEngineMock()
	message := types.Message{
		Payload: map[string]any{
			"action": "single_right",
		},
	}
	args := types.Args{
		"Command":  "$message.action",
		"DeviceId": types.DeviceId("10011cec96"),
	}
	mapping := types.Mapping{
		"Command": {
			"single_left":  "on",
			"single_right": "off",
		},
	}
	err := PostSonoffSwitchMessage([]types.Message{message}, args, mapping, engine)
	s.Nil(err)
}

func (s *ActionsSuite) Test20() {
	engine := mocks.NewEngineMock()
	message := types.Message{}
	args := types.Args{
		"DeviceId": types.DeviceId("0xe0798dfffed39ed1"),
		"Data":     "OFF",
	}
	err := Zigbee2MqttSetState([]types.Message{message}, args, nil, engine)
	s.Nil(err)
	fmt.Println(err)
}

func (s *ActionsSuite) Test30() {
	defer func() { _ = recover() }()
	Get(13)
	s.Fail("expected to panic")
}

func TestActions(t *testing.T) {
	suite.Run(t, new(ActionsSuite))
}
