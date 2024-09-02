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
	// no customization, just call parse_base
	return p.parse_base()
}
