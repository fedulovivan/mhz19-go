package mocks

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var (
	// interface guards
	_ types.DevicesService   = (*mockDevicesService)(nil)
	_ types.ChannelProvider  = (*mockprovider)(nil)
	_ types.EngineAsSupplier = (*mockengine)(nil)
)

type mockDevicesService struct {
}

func (s mockDevicesService) Get() ([]types.Device, error) {
	return nil, nil
}
func (s mockDevicesService) GetByDeviceClass(dc types.DeviceClass) ([]types.Device, error) {
	return nil, nil
}
func (s mockDevicesService) GetOne(id types.DeviceId) (res types.Device, err error) {
	if id == types.DeviceId("10011cec96") {
		res = types.Device{
			Name: "My perfect name",
			Json: map[string]any{
				"Ip":   "192.168.88.60",
				"Port": "8081",
			},
		}
	} else if id == types.DeviceId("nullish-device-id") {
		res = types.Device{}
	} else {
		err = fmt.Errorf("no such device")
	}
	return
}
func (s mockDevicesService) UpsertAll(devices []types.Device) error {
	return nil
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

type mockengine struct {
}

func (e *mockengine) SetDevicesService(s types.DevicesService) {
}

func (e *mockengine) DevicesService() types.DevicesService {
	return &mockDevicesService{}
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

func NewEngineMock() types.EngineAsSupplier {
	return &mockengine{}
}
