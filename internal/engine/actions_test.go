package engine

import (
	"fmt"
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/devices"
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

type mockservice struct {
}

func (s mockservice) Get() ([]devices.Device, error) {
	return nil, nil
}
func (s mockservice) GetOne(id types.DeviceId) (res devices.Device, err error) {
	if id == types.DeviceId("10011cec96") {
		res = devices.Device{
			Json: map[string]any{
				"ip":   "192.168.88.60",
				"port": "8081",
			},
		}
	} else {
		err = fmt.Errorf("no such device")
	}
	return
}

func (s mockservice) Upsert(devices []devices.Device) error {
	return nil
}

func (s *ActionsSuite) Test10() {
	opts := NewOptions()
	opts.SetDevicesService(&mockservice{})
	engine := NewEngine(opts)
	message := types.Message{
		Payload: map[string]any{
			"action": "single_right",
		},
	}
	action := Action{
		Args: Args{
			"Command":        "$message.action",
			"types.DeviceId": types.DeviceId("10011cec96"),
		},
		Mapping: Mapping{
			"Command": {
				"single_left":  "on",
				"single_right": "off",
			},
		},
	}
	PostSonoffSwitchMessage([]types.Message{message}, action, &engine)
}

func TestActions(t *testing.T) {
	suite.Run(t, new(ActionsSuite))
}
