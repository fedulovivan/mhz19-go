package actions

import (
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/mocks"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type ActionsSuite struct {
	suite.Suite
	tag logger.Tag
}

func (s *ActionsSuite) SetupSuite() {
}

func (s *ActionsSuite) TeardownSuite() {
}

func (s *ActionsSuite) Test10() {

	// using this for design tests only
	s.T().Skip()

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
	err := PostSonoffSwitchMessage([]types.Message{message}, args, mapping, engine, s.tag)
	s.Nil(err)
}

func (s *ActionsSuite) Test20() {
	engine := mocks.NewEngineMock()
	message := types.Message{}
	args := types.Args{
		"DeviceId": types.DeviceId("0xe0798dfffed39ed1"),
		"State":    "OFF",
	}
	err := Zigbee2MqttSetState([]types.Message{message}, args, nil, engine, s.tag)
	s.Nil(err)
	fmt.Println(err)
}

func (s *ActionsSuite) Test30() {
	s.PanicsWithValue("Action function 13 not yet implemented", func() {
		Get(13)
	})
}

func (s *ActionsSuite) Test40() {
	engine := mocks.NewEngineMock()
	messages := []types.Message{{}}
	mapping := types.Mapping{}
	args := types.Args{}
	err := YeelightDeviceSetPower(messages, args, mapping, engine, s.tag)
	s.EqualError(err, "no such argument IP")
}

func (s *ActionsSuite) Test41() {
	engine := mocks.NewEngineMock()
	messages := []types.Message{{}}
	mapping := types.Mapping{}
	args := types.Args{
		"IP": "1.1.1.1",
	}
	err := YeelightDeviceSetPower(messages, args, mapping, engine, s.tag)
	s.EqualError(err, "no such argument Cmd")
}

func (s *ActionsSuite) Test42() {
	engine := mocks.NewEngineMock()
	messages := []types.Message{{}}
	mapping := types.Mapping{}
	args := types.Args{
		"IP":  "1.1.1.1",
		"Cmd": "foo",
	}
	err := YeelightDeviceSetPower(messages, args, mapping, engine, s.tag)
	s.EqualError(err, "unsupported command 'foo'")
}

func (s *ActionsSuite) Test43() {

	s.T().Skip()

	engine := mocks.NewEngineMock()
	messages := []types.Message{{}}
	mapping := types.Mapping{}
	args := types.Args{
		"IP":  "192.168.88.169",
		"Cmd": "Off",
	}
	err := YeelightDeviceSetPower(messages, args, mapping, engine, s.tag)
	s.Nil(err)
}

func (s *ActionsSuite) Test50() {

	s.T().Skip()

	mm := []types.Message{{}}
	args := types.Args{}
	mapping := types.Mapping{}
	engine := mocks.NewEngineMock()
	err := PlayAlert(mm, args, mapping, engine, s.tag)
	s.Nil(err)
}

func TestActions(t *testing.T) {
	suite.Run(t, &ActionsSuite{
		tag: logger.NewTag("ActionsSuite"),
	})
}
