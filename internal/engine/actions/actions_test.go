package actions

import (
	"fmt"
	"testing"

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

func (s mockservice) Get() ([]types.Device, error) {
	return nil, nil
}
func (s mockservice) GetByDeviceClass(dc types.DeviceClass) ([]types.Device, error) {
	return nil, nil
}
func (s mockservice) GetOne(id types.DeviceId) (res types.Device, err error) {
	if id == types.DeviceId("10011cec96") {
		res = types.Device{
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
func (s mockservice) UpsertAll(devices []types.Device) error {
	return nil
}

type mockengine struct {
	devicesService types.DevicesService
}

func (e *mockengine) SetDevicesService(s types.DevicesService) {
	e.devicesService = s
}

func (e *mockengine) DevicesService() types.DevicesService {
	return e.devicesService
}

func (e *mockengine) SetMessagesService(s types.MessagesService) {

}

func (e *mockengine) MessagesService() types.MessagesService {
	return nil
}

func (e *mockengine) SetProviders(s ...types.ChannelProvider) {

}

func (e *mockengine) Provider(ct types.ChannelType) types.ChannelProvider {
	return nil
}

func (s *ActionsSuite) Test10() {
	engine := &mockengine{}
	engine.SetDevicesService(&mockservice{})
	message := types.Message{
		Payload: map[string]any{
			"action": "single_right",
		},
	}
	action := types.Action{
		Args: types.Args{
			"Command":  "$message.action",
			"DeviceId": types.DeviceId("10011cec96"),
		},
		Mapping: types.Mapping{
			"Command": {
				"single_left":  "on",
				"single_right": "off",
			},
		},
	}
	s.Nil(PostSonoffSwitchMessage([]types.Message{message}, action, engine))
}

func TestActions(t *testing.T) {
	suite.Run(t, new(ActionsSuite))
}
