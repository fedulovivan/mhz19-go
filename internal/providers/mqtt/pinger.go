package mqtt_provider

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type devicePinger struct {
	parserBase
}

func NewDevicePinger(m MqttLib.Message) *devicePinger {
	return &devicePinger{parserBase{m, types.DEVICE_CLASS_PINGER}}
}

func (p *devicePinger) Parse() (types.Message, bool) {
	out, ok := p.parse_base()
	if ok {
		out.FromEndDevice = true
	}
	return out, ok
}
