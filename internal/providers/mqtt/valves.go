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
	out, ok := p.parse_base()
	if ok {
		out.FromEndDevice = true
	}
	return out, ok
}
