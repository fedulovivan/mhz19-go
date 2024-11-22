package mqtt_provider

import (
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type messagemock struct {
	topic   string
	payload string
}

func (m *messagemock) Topic() string {
	return m.topic
}
func (m *messagemock) Payload() []byte {
	return []byte(m.payload)
}

type ParseSuite struct {
	suite.Suite
}

func (s *ParseSuite) SetupSuite() {

}

func (s *ParseSuite) TeardownSuite() {

}

func (s *ParseSuite) Test10() {
	in := &messagemock{
		topic:   "zigbee2mqtt/0x00158d000405811b",
		payload: `{"power_outage_count":32}`,
	}
	out := Parse(in, types.DEVICE_CLASS_ZIGBEE_DEVICE, true, 1)
	s.Equal(out.ChannelType, types.CHANNEL_MQTT)
	s.Nil(out.RawPayload)
	s.False(out.Timestamp.IsZero())
	s.Equal(out.DeviceClass, types.DEVICE_CLASS_ZIGBEE_DEVICE)
	s.True(out.FromEndDevice)
	s.Equal(out.DeviceId, types.DeviceId("0x00158d000405811b"))
	s.Equal(out.Payload.(map[string]any)["power_outage_count"], float64(32))
	utils.Dump("out", out)
}

func (s *ParseSuite) Test20() {
	in := &messagemock{
		topic:   "device-pinger/192.168.88.1/status",
		payload: `{"status":1}`,
	}
	out := Parse(in, types.DEVICE_CLASS_PINGER, true, 1)
	s.Equal(out.ChannelType, types.CHANNEL_MQTT)
	s.Nil(out.RawPayload)
	s.False(out.Timestamp.IsZero())
	s.Equal(out.DeviceClass, types.DEVICE_CLASS_PINGER)
	s.True(out.FromEndDevice)
	s.Equal(out.DeviceId, types.DeviceId("192.168.88.1"))
	s.Equal(out.Payload.(map[string]any)["status"], float64(1))
	utils.Dump("out", out)
}

func (s *ParseSuite) Test30() {
	in := &messagemock{
		topic:   "zigbee2mqtt/bridge/devices",
		payload: `[{"type": "EndDevice"}]`,
	}
	out := Parse(in, types.DEVICE_CLASS_ZIGBEE_BRIDGE, false, -1)
	s.Equal(out.ChannelType, types.CHANNEL_MQTT)
	s.Nil(out.RawPayload)
	s.False(out.Timestamp.IsZero())
	s.Equal(out.DeviceClass, types.DEVICE_CLASS_ZIGBEE_BRIDGE)
	s.False(out.FromEndDevice)
	s.Zero(out.DeviceId)
	s.Equal(out.Payload.([]any)[0].(map[string]any)["type"], "EndDevice")
	utils.Dump("out", out)
}

func (s *ParseSuite) Test40() {
	in := &messagemock{
		topic:   "espresense/devices/phone:iphone-15/b07cc8",
		payload: `{"distance":4.82}`,
	}
	out := Parse(in, types.DEVICE_CLASS_ESPRESENCE_DEVICE, true, 2)
	s.Equal(out.ChannelType, types.CHANNEL_MQTT)
	s.Nil(out.RawPayload)
	s.False(out.Timestamp.IsZero())
	s.Equal(out.DeviceClass, types.DEVICE_CLASS_ESPRESENCE_DEVICE)
	s.True(out.FromEndDevice)
	s.Equal(out.DeviceId, types.DeviceId("phone:iphone-15"))
	s.Equal(out.Payload.(map[string]any)["distance"], float64(4.82))
	utils.Dump("out", out)
}

func TestArguments(t *testing.T) {
	suite.Run(t, new(ParseSuite))
}
