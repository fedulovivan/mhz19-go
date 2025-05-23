package mocks

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var (
	// interface guards
	_ types.DevicesService  = (*mockDevicesService)(nil)
	_ types.ChannelProvider = (*mockprovider)(nil)
	_ types.ServiceSupplier = (*mockengine)(nil)
)

type mockDevicesService struct {
}

func (s mockDevicesService) Delete(int64) error {
	return nil
}
func (s mockDevicesService) UpdateName(types.Device) error {
	return nil
}
func (s mockDevicesService) UpdateBuriedTimeout(types.Device) error {
	return nil
}
func (s mockDevicesService) Get() ([]types.Device, error) {
	return nil, nil
}
func (s mockDevicesService) GetByDeviceClass(types.DeviceClass) ([]types.Device, error) {
	return nil, nil
}
func (s mockDevicesService) GetOne(id types.DeviceId) (res types.Device, err error) {
	if id == types.DeviceId("10011cec96") {
		name := "My perfect name"
		res = types.Device{
			Name: &name,
			Json: map[string]any{
				"Ip":   "192.168.88.60",
				"Port": "8081",
			},
		}
	} else if id == types.DeviceId("0xe0798dfffed39ed1") {
		name := "My perfect zigbee device"
		res = types.Device{
			Name:        &name,
			DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE,
		}
	} else if id == types.DeviceId("nullish-device-id") {
		res = types.Device{}
	} else {
		err = fmt.Errorf("no such device")
	}
	return
}
func (s mockDevicesService) UpsertAll(devices []types.Device) (int64, error) {
	return 0, nil
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
func (p *mockprovider) Push(m types.Message) {

}

func (p *mockprovider) Channel() types.ChannelType {
	return 0
}

func (p *mockprovider) Type() types.ProviderType {
	return ""
}

func (p *mockprovider) Init() {
}

func (p *mockprovider) Stop() {
}

type mockengine struct {
}

func (e *mockengine) GetDevicesService() types.DevicesService {
	return &mockDevicesService{}
}

func (e *mockengine) GetLdmService() types.LdmService {
	return nil
}

func (e *mockengine) GetMessagesService() types.MessagesService {
	return nil
}

func (e *mockengine) GetProvider(pt types.ProviderType) types.ChannelProvider {
	return &mockprovider{}
}

func NewEngineMock() *mockengine {
	return &mockengine{}
}
