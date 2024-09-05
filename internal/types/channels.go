package types

import "fmt"

type ChannelType byte

const (
	CHANNEL_UNKNOWN  ChannelType = 0
	CHANNEL_MQTT     ChannelType = 1
	CHANNEL_TELEGRAM ChannelType = 2
	CHANNEL_DNS_SD   ChannelType = 3
)

var CHANNEL_NAMES = map[ChannelType]string{
	CHANNEL_UNKNOWN:  "<unknown>",
	CHANNEL_MQTT:     "mqtt",
	CHANNEL_TELEGRAM: "telegram",
	CHANNEL_DNS_SD:   "dns-sd",
}

func (s ChannelType) String() string {
	return fmt.Sprintf("%v (id=%d)", CHANNEL_NAMES[s], s)
}

type ChannelMeta struct {
	MqttTopic string
}

type ChannelProvider interface {
	// getter for the provider's messages channel used in engine.Start
	MessageChan() MessageChan
	// api to invoke provider outbound action, eg:
	// - call tgbotapi.NewMessage for telegram bot provider
	// - post to mqtt topic for mqtt provider
	// - call sonoff http api
	Send(...any) error
	// TODO, api for the unit tests
	Write(m Message)
	Channel() ChannelType
	Init()
	Stop()
}
