package buried_devices_provider

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type provider struct {
	engine.ProviderBase
	ldmRepository ldm.LdmRepository
}

func NewProvider(ldmRepository ldm.LdmRepository) types.ChannelProvider {
	return &provider{
		ldmRepository: ldmRepository,
	}
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_SYSTEM
}

func (p *provider) Init() {
	p.Out = make(types.MessageChan, 100)
	go func() {
		p.ldmRepository.AppendBuriedBlacklist(
			p.ldmRepository.NewKey(types.DEVICE_CLASS_SYSTEM, types.BuriedDeviceId),
			p.ldmRepository.NewKey(types.DEVICE_CLASS_ZIGBEE_BRIDGE, types.DeviceId("bridge")),
		)
		for key := range p.ldmRepository.Buried() {
			time.Sleep(time.Second)
			p.Out <- types.Message{
				ChannelType: types.CHANNEL_SYSTEM,
				DeviceClass: types.DEVICE_CLASS_SYSTEM,
				DeviceId:    types.BuriedDeviceId,
				Timestamp:   time.Now(),
				Payload: map[string]any{
					"BuriedDeviceId": key.DeviceId,
				},
			}
		}
	}()
}
