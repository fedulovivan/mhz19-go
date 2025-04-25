package mqtt_provider

import (
	"encoding/json"
	"log/slog"
	"strings"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

// reduced MqttLib.Message interface
type MqttMessageReduced interface {
	Topic() string
	Payload() []byte
}

// Parse MqttLib.Message into types.Message
// user is oblidged to specify following params:
// deviceClass - DeviceClass
// fromEndDevice - flag indicating whethe this is message from end device or something like z2m bridge or espresense node message
// topicDeviceIdIndex - zero-based index, within topic splitted by slash, where device id part is is located:
// eg for "zigbee2mqtt/0x00158d000405811b" its 1
// for "espresense/devices/apple:10-111/anyvalue" its 2
func Parse(
	msg MqttMessageReduced,
	deviceClass types.DeviceClass,
	fromEndDevice bool,
	topicDeviceIdIndex int,
) types.Message {

	payload := msg.Payload()
	topic := msg.Topic()

	result := types.Message{
		Id:            types.MessageIdSeq.Add(1),
		Timestamp:     time.Now(),
		ChannelType:   types.CHANNEL_MQTT,
		DeviceClass:   deviceClass,
		FromEndDevice: fromEndDevice,
		// DeviceId
		// Payload
		// RawPayload
	}

	tt := strings.Split(strings.TrimLeft(topic, "/"), "/")

	if topicDeviceIdIndex >= 0 && topicDeviceIdIndex < len(tt) {
		result.DeviceId = types.DeviceId(tt[topicDeviceIdIndex])
	}

	if err := json.Unmarshal(payload, &result.Payload); err != nil {
		slog.Warn(tag.F("Failed to parse payload as json"), "payload", string(payload[:]), "err", err)
		result.RawPayload = payload
	}

	return result
}
