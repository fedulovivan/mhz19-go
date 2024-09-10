package types

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

type MessageSuite struct {
	suite.Suite
}

func (s *MessageSuite) SetupSuite() {
}

func (s *MessageSuite) TeardownSuite() {
}

func (s *MessageSuite) Test10() {
	defer func() { _ = recover() }()
	m := Message{}
	v, _ := m.ExecDirective("foo")
	s.Nil(v)
	s.Fail("expected to panic")
}

func (s *MessageSuite) Test20() {

	m := Message{}

	v, err := m.ExecDirective("$channelType")
	s.Equal(CHANNEL_UNKNOWN, v)
	s.Nil(err)

	v, err = m.ExecDirective("$deviceClass")
	s.Equal(DEVICE_CLASS_UNKNOWN, v)
	s.Nil(err)

	v, err = m.ExecDirective("$deviceId")
	s.Equal(DeviceId(""), v)
	s.Nil(err)

}

func (s *MessageSuite) Test30() {

	m := Message{
		ChannelMeta: ChannelMeta{MqttTopic: "foo/111"},
		ChannelType: CHANNEL_MQTT,
		DeviceClass: DEVICE_CLASS_ZIGBEE_DEVICE,
		DeviceId:    "0x00158d0004244bda",
		Payload: map[string]any{
			"action":  "single_left",
			"voltage": 3.33,
			"offline": false,
		},
	}

	v, _ := m.ExecDirective("$message.foo")
	s.Nil(v)

	v, _ = m.ExecDirective("$message.action")
	s.Equal("single_left", v)

	v, _ = m.ExecDirective("$message.voltage")
	s.Equal(3.33, v)

	v, _ = m.ExecDirective("$message.offline")
	s.Equal(false, v)

	v, _ = m.ExecDirective("$channelType")
	s.Equal(CHANNEL_MQTT, v)

	v, _ = m.ExecDirective("$deviceClass")
	s.Equal(DEVICE_CLASS_ZIGBEE_DEVICE, v)

	v, _ = m.ExecDirective("$deviceId")
	s.Equal(DeviceId("0x00158d0004244bda"), v)

}

func (s *MessageSuite) Test40() {
	s.True(IsSpecialDirective("$deviceId"))
}

func (s *MessageSuite) Test50() {
	m := Message{}
	v, _ := m.ExecDirective("$message.foo")
	s.Nil(v)
}
func (s *MessageSuite) Test60() {
	m := Message{
		Payload: map[string]string{
			"foo": "bar",
		},
	}
	v, err := m.ExecDirective("$message.foo")
	s.EqualError(err, `Message.ExecDirective(): Payload is expected to be map[string]any not 'map[string]string', reading field 'foo'`)
	s.Nil(v)
}

func (s *MessageSuite) Test61() {
	m := Message{
		Payload: 111,
	}
	v, err := m.ExecDirective("$message.bar")
	s.EqualError(err, `Message.ExecDirective(): Payload is expected to be map[string]any not 'int', reading field 'bar'`)
	s.Nil(v)
}

func TestMessage(t *testing.T) {
	suite.Run(t, new(MessageSuite))
}
