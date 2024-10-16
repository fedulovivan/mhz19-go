package ldm

import "github.com/fedulovivan/mhz19-go/internal/types"

var _ types.LdmService = (*ldmService)(nil)

type ldmService struct {
	repository LdmRepository
}

func (s ldmService) NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey {
	return s.repository.NewKey(deviceClass, deviceId)
}

func (s ldmService) Get(key types.LdmKey) types.Message {
	return s.repository.Get(key)
}

func (s ldmService) Has(key types.LdmKey) bool {
	return s.repository.Has(key)
}

func (s ldmService) Set(key types.LdmKey, m types.Message) {
	s.repository.Set(key, m)
}

func (s ldmService) GetAll() []types.Message {
	return s.repository.GetAll()
}

func (s ldmService) GetByDeviceId(deviceId types.DeviceId) (types.Message, error) {
	return s.repository.GetByDeviceId(deviceId)
}

func (s ldmService) OnSet() <-chan types.LdmKey {
	return s.repository.OnSet()
}

func NewService(r LdmRepository) types.LdmService {
	return ldmService{
		repository: r,
	}
}
