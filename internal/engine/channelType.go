package engine

import "fmt"

type ChannelType byte

const (
	CHANNEL_UNKNOWN  ChannelType = 0
	CHANNEL_MQTT     ChannelType = 1
	CHANNEL_TELEGRAM ChannelType = 2
)

var CHANNEL_NAMES = map[ChannelType]string{
	CHANNEL_UNKNOWN:  "<unknown>",
	CHANNEL_MQTT:     "mqtt",
	CHANNEL_TELEGRAM: "telegram",
}

func (s ChannelType) String() string {
	return fmt.Sprintf("%v (id=%d)", CHANNEL_NAMES[s], s)
}
