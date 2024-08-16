package mqtt_service

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/engine"
)

type valveManipulator struct {
	parserBase
}

func NewValveManipulator(m MqttLib.Message) *valveManipulator {
	return &valveManipulator{parserBase{m, engine.DEVICE_CLASS_VALVE}}
}

func (p *valveManipulator) Parse() (engine.Message, bool) {
	// no customization, just call parse_base
	return p.parse_base()
}
