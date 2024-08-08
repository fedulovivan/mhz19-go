package mqtt

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/engine"
)

type zigbeeBridge struct {
	parserBase
}

func NewZigbeeBridge(m MqttLib.Message) *zigbeeBridge {
	return &zigbeeBridge{parserBase{m, engine.DEVICE_CLASS_ZIGBEE_BRIDGE}}
}

func (p *zigbeeBridge) Parse() (engine.Message, bool) {
	// no customization, just call parse_base
	return p.parse_base()
}
