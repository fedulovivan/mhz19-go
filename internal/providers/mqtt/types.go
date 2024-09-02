package mqtt_provider

import (
	MqttLib "github.com/eclipse/paho.mqtt.golang"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type TopicHandlers map[string]MqttLib.MessageHandler

type Parser interface {
	Parse() (types.Message, bool)
}

type parserBase struct {
	m  MqttLib.Message
	dc types.DeviceClass
}
