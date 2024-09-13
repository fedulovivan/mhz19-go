package buried_devices_provider

import (
	"fmt"
	"log/slog"
	"sync"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var logTag = logger.MakeTag(logger.BURIED)

type provider struct {
	engine.ProviderBase
	ldmService     types.LdmService
	devicesService types.DevicesService
	buriedTimers   map[types.LdmKey]*time.Timer
	timersMu       sync.Mutex
}

func NewProvider(
	ldmService types.LdmService,
	devicesService types.DevicesService,
) types.ChannelProvider {
	return &provider{
		ldmService:     ldmService,
		devicesService: devicesService,
		buriedTimers:   make(map[types.LdmKey]*time.Timer),
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
			slog.Debug(logTag(fmt.Sprintf("%v device skipped (devices.buried_timeout == 0)", key.DeviceId)))
			skipped = true
		} else {
			slog.Warn(logTag(fmt.Sprintf("%v using custom BuriedTimeout value=%s", key.DeviceId, device.BuriedTimeout.Duration)))
			timeout = device.BuriedTimeout.Duration
		}
	}
	if timer, ok := p.buriedTimers[key]; ok {
		if skipped {
			slog.Warn(logTag(fmt.Sprintf("%v is now skipped, stopping and deleting timer", key.DeviceId)))
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
				p.Out <- types.Message{
					ChannelType: types.CHANNEL_SYSTEM,
					DeviceClass: types.DEVICE_CLASS_SYSTEM,
					DeviceId:    types.DeviceIdForTheBuriedDeviceMessage,
					Timestamp:   time.Now(),
					Payload: map[string]any{
						"BuriedDeviceId": key.DeviceId,
					},
				}
			},
		)
	}
}

func (p *provider) Init() {
	p.Out = make(types.MessageChan, 100)
	go func() {
		for key := range p.ldmService.OnSet() {
			p.handleKey(key)
		}
	}()
}
