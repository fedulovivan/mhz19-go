package types

import (
	"encoding/json"
	"fmt"
)

type ChannelType byte

const (
	CHANNEL_MQTT     ChannelType = 1
	CHANNEL_TELEGRAM ChannelType = 2
	CHANNEL_DNS_SD   ChannelType = 3
	CHANNEL_SYSTEM   ChannelType = 4
	CHANNEL_REST     ChannelType = 5
)

var CHANNEL_NAMES = map[ChannelType]string{
	CHANNEL_MQTT:     "mqtt",
	CHANNEL_TELEGRAM: "telegram",
	CHANNEL_DNS_SD:   "dns-sd",
	CHANNEL_SYSTEM:   "system",
	CHANNEL_REST:     "rest",
}

func (s ChannelType) String() string {
	return CHANNEL_NAMES[s]
}

func (s ChannelType) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%v"`, CHANNEL_NAMES[s])), nil
}

func (s *ChannelType) UnmarshalJSON(b []byte) (err error) {
	var v string
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	for channel, name := range CHANNEL_NAMES {
		if fmt.Sprintf("ChannelType(%s)", name) == v || fmt.Sprintf("ChannelType(%d)", channel) == v {
			*s = channel
			return
		}
	}
	return fmt.Errorf("failed to unmarshal %v (type=%T) to ChannelType", v, v)
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

	Type() ProviderType

	// a channel type this provider was created for
	Channel() ChannelType

	// init provider
	Init()

	// stop provider
	Stop()

	// api to push message to channel
	Push(m Message)
}
