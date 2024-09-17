package mqtt_provider

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type zigbeeDevice struct {
	parserBase
}

func NewZigbeeDevice(m MqttLib.Message) *zigbeeDevice {
	return &zigbeeDevice{parserBase{m, types.DEVICE_CLASS_ZIGBEE_DEVICE}}
}

func (p *zigbeeDevice) Parse() (types.Message, bool) {
	out, ok := p.parse_base()
	if ok {
		out.FromEndDevice = true
	}
	return out, ok
}
