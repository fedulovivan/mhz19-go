package last_device_message

import "github.com/fedulovivan/mhz19-go/internal/types"

type ldmService struct {
	repository LdmRepository
}

func (s ldmService) MakeKey(deviceClass types.DeviceClass, deviceId types.DeviceId) types.LdmKey {
	return s.repository.MakeKey(deviceClass, deviceId)
}

func (s ldmService) Get(key types.LdmKey) types.Message {
	return s.repository.Get(key)
}

func (s ldmService) Set(key types.LdmKey, m types.Message) {
	s.repository.Set(key, m)
}

func (s ldmService) GetAll() []types.Message {
	return s.repository.GetAll()
}

func (s ldmService) GetByDeviceId(deviceId types.DeviceId) types.Message {
	return s.repository.GetByDeviceId(deviceId)
}

func NewService(r LdmRepository) types.LdmService {
	return ldmService{
		repository: r,
	}
}
