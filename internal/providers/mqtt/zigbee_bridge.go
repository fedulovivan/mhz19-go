package mqtt_provider

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type zigbeeBridge struct {
	parserBase
}

func NewZigbeeBridge(m MqttLib.Message) *zigbeeBridge {
	return &zigbeeBridge{parserBase{m, types.DEVICE_CLASS_ZIGBEE_BRIDGE}}
}

func (p *zigbeeBridge) Parse() (types.Message, bool) {
	// no customization, just call parse_base
	return p.parse_base()
}
