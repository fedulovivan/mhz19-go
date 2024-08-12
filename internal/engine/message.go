package engine

import (
	"fmt"
	"time"
)

type ChannelMeta struct {
	MqttTopic string
}

type Message struct {
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
	RawPayload []byte
	// additional metadata specific for the current channel
	ChannelMeta ChannelMeta
}

// tuple of current and previous messages
type MessageTuple = [2]Message

// get message primitive field or message payload field
func (m *Message) Get(field string) (any, error) {
	switch field {
	case "ChannelMeta":
		return m.ChannelMeta, nil
	case "ChannelType":
		return m.ChannelType, nil
	case "DeviceClass":
		return m.DeviceClass, nil
	case "DeviceId":
		return m.DeviceId, nil
	case "Timestamp":
		return m.Timestamp, nil
	default:
		p, ok := m.Payload.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Message.Payload is not a map[string]any, reading field [%v]", field)
		}
		v, ok := p[field]
		if !ok {
			return nil, fmt.Errorf("Message.Payload has no field [%v]", field)
		}
		return v, nil
	}
	// fmt.Printf("%+v, %T", m.Payload, m.Payload)
	// case "Payload":
	// 	return m.Payload, nil
	// case "RawPayload":
	// 	return m.RawPayload, nil
}

type MessageChan chan Message
