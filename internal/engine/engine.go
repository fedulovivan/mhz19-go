package engine

import (
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/logger"
)

var withTag = logger.MakeTag("ENGN")

var services []Service

func Start(input ...Service) {
	services = input
	start()
}

func start() {
	for _, service := range services {
		go func(s Service) {
			s.Init()
			for m := range s.Receive() {
				HandleMessage(m)
			}
		}(service)
	}
}

func Stop() {
	for _, s := range services {
		s.Stop()
	}
}

func HandleMessage(m Message) {

	p := m.Payload
	if m.DeviceClass == DEVICE_CLASS_ZIGBEE_BRIDGE {
		p = "<too big to render>"
	}

	slog.Debug(
		withTag("New message"),
		"ChannelType", m.ChannelType,
		"ChannelMeta", m.ChannelMeta,
		"DeviceClass", m.DeviceClass,
		"DeviceId", m.DeviceId,
		"Payload", p,
	)
}
