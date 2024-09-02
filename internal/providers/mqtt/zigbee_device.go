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
	// no customization, just call parse_base
	return p.parse_base()
}
