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

type TimerWithLastSeen struct {
	lastSeen time.Time
	timer    *time.Timer
}

type Timers map[types.LdmKey]*TimerWithLastSeen

type provider struct {
	engine.ProviderBase
	ldmService     types.LdmService
	devicesService types.DevicesService
	buriedTimers   Timers
	timersMu       sync.RWMutex
}

var _ types.ChannelProvider = (*provider)(nil)

func NewProvider(
	ldmService types.LdmService,
	devicesService types.DevicesService,
) *provider {
	return &provider{
		ldmService:     ldmService,
		devicesService: devicesService,
		buriedTimers:   make(Timers),
	}
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_SYSTEM
}

func (p *provider) handleKey(key types.LdmKey) {
	p.timersMu.Lock()
	defer p.timersMu.Unlock()
	skipped := false
	// device whether we need to skip emitting message for certain device
	// or we need to use custom timeout for this device
	timeout := app.Config.DefaultBuriedTimeout
	device, err := p.devicesService.GetOne(key.DeviceId)
	if err == nil && device.BuriedTimeout != nil {
		if device.BuriedTimeout.Duration == 0 {
			slog.Debug(tag.F("%v device skipped (devices.buried_timeout == 0)", key.DeviceId))
			skipped = true
		} else {
			slog.Debug(tag.F("%v using custom BuriedTimeout value=%s", key.DeviceId, device.BuriedTimeout.Duration))
			timeout = device.BuriedTimeout.Duration
		}
	}
	// timerWithLastSeen already exist
	if timerWithLastSeen, ok := p.buriedTimers[key]; ok {
		// delete timer if on next reading of device data "skipped" flag has been changed to true
		if skipped {
			slog.Warn(tag.F("%v is now skipped, stopping and deleting timer", key.DeviceId))
			timerWithLastSeen.timer.Stop()
			delete(p.buriedTimers, key)
			return
		}
		// prolong timer after receiving next message
		// and also check if timer was already fired and emit approprite message
		active := timerWithLastSeen.timer.Reset(timeout)
		if !active {
			p.emitMessage(
				key.DeviceId,
				"ceased",
				timerWithLastSeen.lastSeen, // read before updating, to celculated duration of missing period for "ceased" message
			)
		}
		timerWithLastSeen.lastSeen = time.Now()
	} else {
		if skipped {
			return
		}
		// register timer
		p.buriedTimers[key] = &TimerWithLastSeen{
			timer: time.AfterFunc(
				timeout,
				func() {
					p.timersMu.RLock()
					lastSeen := p.buriedTimers[key].lastSeen
					p.timersMu.RUnlock()
					p.emitMessage(
						key.DeviceId,
						"fired",
						lastSeen,
					)
				},
			),
			lastSeen: time.Now(),
		}
	}
}

func (p *provider) emitMessage(
	deviceId types.DeviceId,
	transition string,
	lastSeen time.Time,
) {
	outMsg := types.Message{
		Id:          types.MessageIdSeq.Add(1),
		Timestamp:   time.Now(),
		ChannelType: types.CHANNEL_SYSTEM,
		DeviceClass: types.DEVICE_CLASS_SYSTEM,
		DeviceId:    types.DEVICE_ID_FOR_THE_BURIED_DEVICES_PROVIDER_MESSAGE,
		Payload: map[string]any{
			"BuriedDeviceId": deviceId,
			"Transition":     transition,
			"HaveNotSeen":    time.Since(lastSeen).Round(time.Second),
		},
		FromEndDevice: false,
	}
	p.Push(outMsg)
}

func (p *provider) Init() {
	p.ProviderBase.Init()
	go func() {
		for key := range p.ldmService.OnSet() {
			p.handleKey(key)
		}
	}()
}
