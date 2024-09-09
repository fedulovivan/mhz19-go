package types

import (
	"fmt"
	"strings"
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
	// message payload json specific for channel and device class
	Payload any `json:"payload,omitempty"`
	// filled only if failed to parse into json
	RawPayload []byte `json:"-"`
	// additional metadata specific for the current channel
	ChannelMeta ChannelMeta `json:"-"`
}

type MessageChan chan Message

// tuple of current and previous messages, prev could be nil
type MessageTuple struct {
	Curr *Message
	Prev *Message
}

func (m *Message) IsSpecial(field string) bool {
	return field == "$deviceId" || field == "$deviceClass" || field == "$channelType" || strings.HasPrefix(field, "$message.")
}

// read message or message payload field using special syntax designed to be used in types.Args
func (m *Message) ExecDirective(field string) (any, error) {
	if field == "$deviceId" {
		return m.DeviceId, nil
	} else if field == "$deviceClass" {
		return m.DeviceClass, nil
	} else if field == "$channelType" {
		return m.ChannelType, nil
	} else if strings.HasPrefix(field, "$message.") {
		_, field, _ := strings.Cut(field, ".")
		if m.Payload == nil {
			return nil, nil
		}
		p, ok := m.Payload.(map[string]any)
		if !ok {
			return nil, fmt.Errorf("Message.ExecDirective(): Payload is expected to be map[string]any not '%T', reading field '%v'", m.Payload, field)
		}
		v, ok := p[field]
		if !ok {
			return nil, fmt.Errorf("Message.ExecDirective(): Payload '%T, %+v' has no field '%v'", m.Payload, m.Payload, field)
		}
		return v, nil
	}
	return nil, nil
}
