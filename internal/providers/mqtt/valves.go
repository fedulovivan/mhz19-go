package mqtt_provider

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type valveManipulator struct {
	parserBase
}

func NewValveManipulator(m MqttLib.Message) *valveManipulator {
	return &valveManipulator{parserBase{m, types.DEVICE_CLASS_VALVE}}
}

func (p *valveManipulator) Parse() (types.Message, bool) {
	// no customization, just call parse_base
	return p.parse_base()
}
