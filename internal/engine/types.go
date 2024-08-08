package engine

import (
	"fmt"
	"time"
)

type JsonPayload interface{}

// type JsonPayload map[string]interface{}

type DeviceId string

type DeviceClass byte

const (
	DEVICE_CLASS_ZIGBEE_DEVICE DeviceClass = 1
	DEVICE_CLASS_PINGER        DeviceClass = 2
	DEVICE_CLASS_VALVE         DeviceClass = 3
	DEVICE_CLASS_ZIGBEE_BRIDGE DeviceClass = 4
)

var DEVICE_CLASS_NAMES = map[DeviceClass]string{
	DEVICE_CLASS_ZIGBEE_DEVICE: "zigbee-device",
	DEVICE_CLASS_PINGER:        "device-pinger",
	DEVICE_CLASS_VALVE:         "valve-manipulator",
	DEVICE_CLASS_ZIGBEE_BRIDGE: "zigbee-bridge",
}

func (s DeviceClass) String() string {
	if s == 0 {
		return "<unknown>"
	}
	return fmt.Sprintf("%v (id=%d)", DEVICE_CLASS_NAMES[s], s)
}

type ChannelType byte

const (
	CHANNEL_MQTT     ChannelType = 1
	CHANNEL_TELEGRAM ChannelType = 2
)

var CHANNEL_NAMES = map[ChannelType]string{
	CHANNEL_MQTT:     "mqtt",
	CHANNEL_TELEGRAM: "telegram",
}

func (s ChannelType) String() string {
	if s == 0 {
		return "<unknown>"
	}
	return fmt.Sprintf("%v (id=%d)", CHANNEL_NAMES[s], s)
}

type ChannelMeta struct {
	MqttTopic string
}

type Message struct {
	// additional metadata specific for the current channel
	ChannelMeta ChannelMeta
	// channel, which was used to receive message
	ChannelType ChannelType `json:"channelType"`
	// device class, see DeviceClass
	DeviceClass DeviceClass `json:"deviceClass"`
	// device id extracted from topic
	DeviceId DeviceId `json:"deviceId"`
	// time when message was received by backend
	Timestamp time.Time `json:"timestamp"`
	// parsed message payload json
	Payload JsonPayload `json:"payload"`
	// filled only if failed to parse into json
	RawPayload []byte `json:"rawPayload"`
}

type MessageChan chan Message

type Service interface {
	Receive() MessageChan
	Type() ChannelType
	Init()
	Stop()
}
