package types

import "fmt"

type ChannelType byte

const (
	CHANNEL_MQTT     ChannelType = 1
	CHANNEL_TELEGRAM ChannelType = 2
	CHANNEL_DNS_SD   ChannelType = 3
	CHANNEL_SYSTEM   ChannelType = 4
)

var CHANNEL_NAMES = map[ChannelType]string{
	CHANNEL_MQTT:     "mqtt",
	CHANNEL_TELEGRAM: "telegram",
	CHANNEL_DNS_SD:   "dns-sd",
	CHANNEL_SYSTEM:   "system",
}

func (s ChannelType) String() string {
	return fmt.Sprintf("%v (id=%d)", CHANNEL_NAMES[s], s)
}

func (s ChannelType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, CHANNEL_NAMES[s])), nil
}

type ChannelMeta struct {
	MqttTopic string
}

type ChannelProvider interface {

	// getter for the provider's messages channel
	Messages() MessageChan

	// api to invoke provider outbound action, eg:
	// - call tgbotapi.NewMessage for telegram bot provider
	// - post to mqtt topic for mqtt provider
	// - call sonoff http api
	Send(...any) error

	// a channel type this provider was created for
	Channel() ChannelType

	Init()
	Stop()
}
