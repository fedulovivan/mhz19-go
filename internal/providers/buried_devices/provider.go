package buried_devices_provider

import (
	"log/slog"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var tag = utils.NewTag(logger.BURIED)

type BuriedTimers map[types.LdmKey]*time.Timer

type provider struct {
	engine.ProviderBase
	ldmService     types.LdmService
	devicesService types.DevicesService
	buriedTimers   BuriedTimers
	timersMu       sync.Mutex
}

var _ types.ChannelProvider = (*provider)(nil)

func NewProvider(
	ldmService types.LdmService,
	devicesService types.DevicesService,
) *provider {
	return &provider{
		ldmService:     ldmService,
		devicesService: devicesService,
		buriedTimers:   make(BuriedTimers),
	}
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_SYSTEM
}

func (p *provider) handleKey(key types.LdmKey) {
	p.timersMu.Lock()
	defer p.timersMu.Unlock()
	skipped := false
	timeout := app.Config.DefaultBuriedTimeout
	device, err := p.devicesService.GetOne(key.DeviceId)
	if err == nil && device.BuriedTimeout != nil {
		if device.BuriedTimeout.Duration == 0 {
			slog.Debug(tag.F("%v device skipped (devices.buried_timeout == 0)", key.DeviceId))
			skipped = true
		} else {
			slog.Warn(tag.F("%v using custom BuriedTimeout value=%s", key.DeviceId, device.BuriedTimeout.Duration))
			timeout = device.BuriedTimeout.Duration
		}
	}
	if timer, ok := p.buriedTimers[key]; ok {
		if skipped {
			slog.Warn(tag.F("%v is now skipped, stopping and deleting timer", key.DeviceId))
			timer.Stop()
			delete(p.buriedTimers, key)
			return
		}
		timer.Reset(timeout)
	} else {
		if skipped {
			return
		}
		p.buriedTimers[key] = time.AfterFunc(
			timeout,
			func() {
				outMsg := types.Message{
					Id:          types.MessageIdSeq.Add(1),
					Timestamp:   time.Now(),
					ChannelType: types.CHANNEL_SYSTEM,
					DeviceClass: types.DEVICE_CLASS_SYSTEM,
					DeviceId:    types.DEVICE_ID_FOR_THE_BURIED_DEVICES_PROVIDER_MESSAGE,
					Payload: map[string]any{
						"BuriedDeviceId": key.DeviceId,
					},
					FromEndDevice: false,
				}
				p.Push(outMsg)
			},
		)
	}
}

func (p *provider) Init() {
	p.ProviderBase.Init()
	go func() {
		for key := range p.ldmService.OnSet() {
			p.handleKey(key)
		}
	}()
}
