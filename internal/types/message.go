package types

import (
	"fmt"
	"strings"
	"time"

	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var MessageIdSeq = utils.NewSeq(0)

type Message struct {
	// unique id
	Id int32 `json:"-"`
	// channel, which was used to receive message
	ChannelType ChannelType `json:"channelType"`
	// device class, see DeviceClass
	DeviceClass DeviceClass `json:"deviceClass"`
	// device id, specific for the current channel and device class, eg ieee address for zigbee device
	DeviceId DeviceId `json:"deviceId"`
	// indicates this is end device message,
	// and not a thing like z2m bridge message with list of registered devices,
	// or not a dns-sd channel message with sonoff device announcement
	// or not system message from buried_devices/provider.go
	FromEndDevice bool `json:"fromEndDevice"`
	// time when message was received by backend
	Timestamp time.Time `json:"timestamp"`
	// message payload json specific for channel and device class
	Payload any `json:"payload,omitempty"`
	// filled only if failed to parse into json
	RawPayload []byte `json:"-"`
	// additional metadata specific for the current channel
	ChannelMeta *ChannelMeta `json:"-"`
}

type MessageChan chan Message

// container to pass current, previous messages and contents of queue messages for throttled rules
// when its normal handling, Curr is always specified, while Prev could be empty
// (if this is first message for such device class and device id and LdmService.Has() returns false)
// however, when Condition.OtherDeviceId is set, only Curr may be filled if exists in LdmService
// (but also could be not filled, meaning both may be nil)
// Throttled - TBD
type MessageCompound struct {
	Curr   *Message
	Prev   *Message
	Queued []Message
}

type GetCompoundForOtherDeviceId = func(otherDeviceId DeviceId) MessageCompound

func IsSpecialDirective(field string) bool {
	return field == "$deviceId" || field == "$deviceClass" || field == "$channelType" || field == "$fromEndDevice" || strings.HasPrefix(field, "$message.")
}

// read message or message payload field using special syntax designed to be used in types.Args
func (m *Message) ExecDirective(field string) (any, error) {
	if m == nil {
		return nil, nil
		// panic("Message.ExecDirective(): message is nil")
	}
	if field == "$deviceId" {
		return m.DeviceId, nil
	} else if field == "$deviceClass" {
		return m.DeviceClass, nil
	} else if field == "$channelType" {
		return m.ChannelType, nil
	} else if field == "$fromEndDevice" {
		return m.FromEndDevice, nil
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
			// return nil, fmt.Errorf("Message.ExecDirective(): Payload '%T, %+v' has no field '%v'", m.Payload, m.Payload, field)
			return nil, nil
		}
		return v, nil
	} else {
		panic(fmt.Sprintf("unknown directive %s", field))
	}
}

func NewSystemMessage(text string) Message {
	return Message{
		Id:            MessageIdSeq.Inc(),
		Timestamp:     time.Now(),
		ChannelType:   CHANNEL_SYSTEM,
		DeviceClass:   DEVICE_CLASS_SYSTEM,
		DeviceId:      DEVICE_ID_FOR_THE_APPLICATION_MESSAGE,
		FromEndDevice: false,
		Payload: map[string]any{
			"text": text,
		},
	}
}

// func NewMessage(
// 	ct ChannelType,
// 	// fromEndDevice bool,
// 	// dc *DeviceClass,
// 	// deviceId *DeviceId,
// ) Message {
// 	res := Message{
// 		Id:          IdSeq.Inc(),
// 		Timestamp:   time.Now(),
// 		ChannelType: ct,
// 	}
// 	// FromEndDevice: fromEndDevice,
// 	// if fromEndDevice && dc == nil && deviceId == nil {
// 	// 	panic("DeviceClass and DeviceId are mandatory for end device message")
// 	// }
// 	// if dc != nil {
// 	// 	res.DeviceClass = *dc
// 	// }
// 	// if deviceId != nil {
// 	// 	res.DeviceId = *deviceId
// 	// }
// 	return res
// }
