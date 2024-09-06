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

type mockprovider struct {
}

func (p *mockprovider) Messages() types.MessageChan {
	return nil
}
func (p *mockprovider) Send(a ...any) error {
	fmt.Println(a...)
	return nil
}
func (p *mockprovider) Write(m types.Message) {

}
func (p *mockprovider) Channel() types.ChannelType {
	return types.CHANNEL_UNKNOWN
}

func (p *mockprovider) Init() {

}
func (p *mockprovider) Stop() {

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
				"Ip":   "192.168.88.60",
				"Port": "8081",
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
}

func (e *mockengine) SetDevicesService(s types.DevicesService) {
}

func (e *mockengine) DevicesService() types.DevicesService {
	return &mockservice{}
}

func (e *mockengine) SetMessagesService(s types.MessagesService) {

}

func (e *mockengine) MessagesService() types.MessagesService {
	return nil
}

func (e *mockengine) SetProviders(s ...types.ChannelProvider) {

}

func (e *mockengine) Provider(ct types.ChannelType) types.ChannelProvider {
	return &mockprovider{}
}

func (s *ActionsSuite) Test10() {
	engine := &mockengine{}
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
	engine := &mockengine{}
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
