package mqtt

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/engine"
)

type zigbeeDevice struct {
	parserBase
}

func NewZigbeeDevice(m MqttLib.Message) *zigbeeDevice {
	return &zigbeeDevice{parserBase{m, engine.DEVICE_CLASS_ZIGBEE_DEVICE}}
}

func (p *zigbeeDevice) Parse() (engine.Message, bool) {
	// no customization, just call parse_base
	return p.parse_base()
}
