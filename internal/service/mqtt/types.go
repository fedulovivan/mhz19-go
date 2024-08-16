package mqtt_service

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/engine"
)

type TopicHandlers map[string]MqttLib.MessageHandler

type Parser interface {
	Parse() (engine.Message, bool)
}

type parserBase struct {
	m  MqttLib.Message
	dc engine.DeviceClass
}
