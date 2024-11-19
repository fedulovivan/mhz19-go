package actions

import (
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/mocks"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type ActionsSuite struct {
	suite.Suite
	tag utils.Tag
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
	compound := types.MessageCompound{
		Curr: &message,
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
	err := PostSonoffSwitchMessage(compound, args, mapping, engine, s.tag)
	s.Nil(err)
}

func (s *ActionsSuite) Test20() {
	engine := mocks.NewEngineMock()
	message := types.Message{}
	compound := types.MessageCompound{
		Curr: &message,
	}

	args := types.Args{
		"DeviceId": types.DeviceId("0xe0798dfffed39ed1"),
		"State":    "OFF",
	}
	err := MqttSetState(compound, args, nil, engine, s.tag)
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
	message := types.Message{}
	compound := types.MessageCompound{
		Curr: &message,
	}
	mapping := types.Mapping{}
	args := types.Args{}
	err := YeelightDeviceSetPower(compound, args, mapping, engine, s.tag)
	s.EqualError(err, "no such argument IP")
}

func (s *ActionsSuite) Test41() {
	engine := mocks.NewEngineMock()
	message := types.Message{}
	compound := types.MessageCompound{
		Curr: &message,
	}
	mapping := types.Mapping{}
	args := types.Args{
		"IP": "1.1.1.1",
	}
	err := YeelightDeviceSetPower(compound, args, mapping, engine, s.tag)
	s.EqualError(err, "no such argument Cmd")
}

func (s *ActionsSuite) Test42() {
	engine := mocks.NewEngineMock()
	message := types.Message{}
	compound := types.MessageCompound{
		Curr: &message,
	}
	mapping := types.Mapping{}
	args := types.Args{
		"IP":  "1.1.1.1",
		"Cmd": "foo",
	}
	err := YeelightDeviceSetPower(compound, args, mapping, engine, s.tag)
	s.EqualError(err, "unsupported command 'foo'")
}

func (s *ActionsSuite) Test43() {

	s.T().Skip()

	engine := mocks.NewEngineMock()
	message := types.Message{}
	compound := types.MessageCompound{
		Curr: &message,
	}
	mapping := types.Mapping{}
	args := types.Args{
		"IP":  "192.168.88.169",
		"Cmd": "Off",
	}
	err := YeelightDeviceSetPower(compound, args, mapping, engine, s.tag)
	s.Nil(err)
}

func (s *ActionsSuite) Test50() {

	s.T().Skip()

	message := types.Message{}
	compound := types.MessageCompound{
		Curr: &message,
	}
	args := types.Args{}
	mapping := types.Mapping{}
	engine := mocks.NewEngineMock()
	err := PlayAlert(compound, args, mapping, engine, s.tag)
	s.Nil(err)
}

func TestActions(t *testing.T) {
	suite.Run(t, &ActionsSuite{
		tag: utils.NewTag("ActionsSuite"),
	})
}
