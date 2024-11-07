package ldm

import "github.com/fedulovivan/mhz19-go/internal/types"

var _ types.LdmService = (*service)(nil)

type service struct {
	repository LdmRepository
}

func (s service) NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey {
	return s.repository.NewKey(deviceClass, deviceId)
}

func (s service) Get(key types.LdmKey) types.Message {
	return s.repository.Get(key)
}

func (s service) Has(key types.LdmKey) bool {
	return s.repository.Has(key)
}

func (s service) Set(key types.LdmKey, m types.Message) {
	s.repository.Set(key, m)
}

func (s service) GetAll() []types.Message {
	return s.repository.GetAll()
}

func (s service) GetByDeviceId(deviceId types.DeviceId) (types.Message, error) {
	return s.repository.GetByDeviceId(deviceId)
}

func (s service) OnSet() <-chan types.LdmKey {
	return s.repository.OnSet()
}

func NewService(r LdmRepository) service {
	return service{
		repository: r,
	}
}
