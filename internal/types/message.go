package types

import (
	"fmt"
	"time"
)

type Message struct {
	// channel, which was used to receive message
	ChannelType ChannelType `json:"channelType"`
	// device class, see DeviceClass
	DeviceClass DeviceClass `json:"deviceClass"`
	// device id, specific for the current channel and device class, eg ieee adress for zigbee device
	DeviceId DeviceId `json:"deviceId"`
	// time when message was received by backend
	Timestamp time.Time `json:"timestamp"`
	// parsed message payload json
	Payload any `json:"payload,omitempty"`
	// filled only if failed to parse into json
	RawPayload []byte `json:"-"`
	// additional metadata specific for the current channel
	ChannelMeta ChannelMeta `json:"-"`
}

// tuple of current and previous messages
type MessageTuple = [2]Message

// get message primitive field or message payload field
func (m *Message) Get(field string) (any, error) {
	switch field {
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
			return nil, fmt.Errorf("Message.Get(): Payload expected to be map[string]any instead of '%T', reading field '%v'", m.Payload, field)
		}
		v, ok := p[field]
		if !ok {
			return nil, fmt.Errorf("Message.Get(): Payload '%T' has no field '%v'", m.Payload, field)
		}
		return v, nil
	}
	// case "ChannelMeta":
	// 	return m.ChannelMeta, nil
	// fmt.Printf("%+v, %T", m.Payload, m.Payload)
	// case "Payload":
	// 	return m.Payload, nil
	// case "RawPayload":
	// 	return m.RawPayload, nil
}

type MessageChan chan Message
