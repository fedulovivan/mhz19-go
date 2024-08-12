package engine

import (
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
}

func (s *EngineSuite) Test10() {
	actual := matchesCondition(MessageTuple{}, Condition{})
	s.False(actual)
}

func (s *EngineSuite) Test11() {
	actual := matchesCondition(MessageTuple{}, Condition{
		Or: true,
		List: []Condition{
			{Fn: Equal, Args: NamedArgs{"Left": true, "Right": false}},
			{Fn: Equal, Args: NamedArgs{"Left": 1.11, "Right": 1.11}},
		},
	})
	s.True(actual)
}

func (s *EngineSuite) Test12() {
	actual := matchesCondition(MessageTuple{}, Condition{
		List: []Condition{
			{Fn: Equal, Args: NamedArgs{"Left": true, "Right": false}},
			{Fn: Equal, Args: NamedArgs{"Left": 1.11, "Right": 1.11}},
		},
	})
	s.False(actual)
}

func (s *EngineSuite) Test20() {
	defer func() { _ = recover() }()
	matchFunction(MessageTuple{}, "foo1", nil)
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test30() {
	actual := listSome(MessageTuple{}, []Condition{})
	s.False(actual)
}

func (s *EngineSuite) Test31() {
	actual := listSome(MessageTuple{}, []Condition{
		{Fn: Equal, Args: NamedArgs{"Left": 1, "Right": 1}},
		{Fn: Equal, Args: NamedArgs{"Left": "foo", "Right": "bar"}},
	})
	s.True(actual)
}

func (s *EngineSuite) Test40() {
	actual := listEvery(MessageTuple{}, []Condition{})
	s.False(actual)
}

func (s *EngineSuite) Test41() {
	actual := listEvery(MessageTuple{}, []Condition{
		{Fn: Equal, Args: NamedArgs{"Left": 1, "Right": 1}},
		{Fn: Equal, Args: NamedArgs{"Left": "foo", "Right": "foo"}},
	})
	s.True(actual)
}

func (s *EngineSuite) Test50() {
	actual, err := arg(Message{}, NamedArgs{"foo": 1}, "foo")
	expected := 1
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test51() {
	actual, err := arg(Message{DeviceId: "foo2"}, NamedArgs{"Lorem": "$message.DeviceId"}, "Lorem")
	expected := DeviceId("foo2")
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test52() {
	actual, err := arg(Message{DeviceId: "foo3"}, NamedArgs{"Lorem1": "$deviceId"}, "Lorem1")
	expected := DeviceId("foo3")
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test53() {
	m := Message{
		Payload: map[string]any{
			"action": "my_action",
		},
	}
	args := NamedArgs{
		"Lorem3": "$message.action",
	}
	actual, err := arg(m, args, "Lorem3")
	expected := "my_action"
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test54() {
	m := Message{
		Payload: map[string]any{
			"voltage": 3.33,
		},
	}
	args := NamedArgs{
		"Lorem3": "$message.voltage",
	}
	actual, err := arg(m, args, "Lorem3")
	expected := 3.33
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test55() {
	m := Message{
		Payload: map[string]any{
			"foo": false,
		},
	}
	args := NamedArgs{
		"Lorem3": "$message.bar",
	}
	actual, err := arg(m, args, "Lorem3")
	s.Nil(actual)
	s.NotEmpty(err)
}

func (s *EngineSuite) Test56() {
	actual, err := arg(Message{DeviceClass: DEVICE_CLASS_ZIGBEE_BRIDGE}, NamedArgs{"Lorem1": "$deviceClass"}, "Lorem1")
	s.Equal(DEVICE_CLASS_ZIGBEE_BRIDGE, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test60() {
	executeActions([]Action{})
}

func (s *EngineSuite) Test70() {
	handleMessage(Message{}, []Rule{})
}

func (s *EngineSuite) Test71() {
	handleMessage(Message{DeviceClass: DEVICE_CLASS_ZIGBEE_BRIDGE}, []Rule{})
}

func (s *EngineSuite) Test72() {
	handleMessage(Message{}, []Rule{
		{
			Condition: Condition{
				Fn:   Equal,
				Args: NamedArgs{"Left": true, "Right": true},
			},
		},
	})
}

func (s *EngineSuite) Test73() {
	defer func() { _ = recover() }()
	handleMessage(Message{}, []Rule{
		{
			Condition: Condition{
				Fn:   Equal,
				Args: NamedArgs{"Left": true, "Right": true},
			},
			Throttle: time.Minute,
		},
	})
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test80() {
	actual := EqualFn(MessageTuple{}, NamedArgs{})
	s.False(actual)
}

func (s *EngineSuite) Test81() {
	actual := EqualFn(MessageTuple{}, NamedArgs{"Left": 1, "Right": 1})
	s.True(actual)
}

func (s *EngineSuite) Test82() {
	actual := EqualFn(MessageTuple{}, NamedArgs{"Left": "one", "Right": "one"})
	s.True(actual)
}

func (s *EngineSuite) Test83() {
	actual := EqualFn(
		MessageTuple{
			Message{
				Payload: map[string]any{
					"action": "my_action",
				},
			},
		},
		NamedArgs{"Left": "$message.action", "Right": "my_action"},
	)
	s.True(actual)
}

func (s *EngineSuite) Test90() {
	actual := NotEqualFn(MessageTuple{}, NamedArgs{})
	s.False(actual)
}

func (s *EngineSuite) Test91() {
	actual := NotEqualFn(MessageTuple{}, NamedArgs{"Left": 1, "Right": 1})
	s.False(actual)
}

func (s *EngineSuite) Test92() {
	actual := NotEqualFn(MessageTuple{}, NamedArgs{"Left": "one", "Right": "one"})
	s.False(actual)
}

func (s *EngineSuite) Test93() {
	mt := MessageTuple{
		Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := NamedArgs{"Left": "$message.action", "Right": "my_action"}
	actual := NotEqualFn(
		mt,
		args,
	)
	s.False(actual)
}

func (s *EngineSuite) Test100() {
	actual := InListFn(MessageTuple{}, NamedArgs{})
	s.False(actual)
}

func (s *EngineSuite) Test101() {
	mt := MessageTuple{
		Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := NamedArgs{
		"Value": "$message.action",
		"List": []any{
			"foo1",
			"my_action",
		},
	}
	actual := InListFn(mt, args)
	s.True(actual)
}

func (s *EngineSuite) Test102() {
	mt := MessageTuple{
		Message{
			Payload: map[string]any{
				"voltage": 1.11,
			},
		},
	}
	args := NamedArgs{
		"Value": "$message.voltage",
		"List": []any{
			1,
			1.11,
			2.0,
		},
	}
	actual := InListFn(mt, args)
	s.True(actual)
}

func (s *EngineSuite) Test103() {
	mt := MessageTuple{
		Message{},
	}
	args := NamedArgs{
		"Value": "some1",
		"List":  []any{"some2", "some3"},
	}
	actual := InListFn(mt, args)
	s.False(actual)
}

func (s *EngineSuite) Test104() {
	args := NamedArgs{
		"Value": "some1",
		"List":  "some2",
	}
	actual := InListFn(MessageTuple{}, args)
	s.False(actual)
}

func (s *EngineSuite) Test110() {
	actual := NotNilFn(MessageTuple{}, NamedArgs{})
	s.False(actual)
}

func (s *EngineSuite) Test111() {
	actual := NotNilFn(MessageTuple{}, NamedArgs{"Value": "foo"})
	s.True(actual)
}

func (s *EngineSuite) Test112() {
	actual := NotNilFn(MessageTuple{}, NamedArgs{"Value": false})
	s.True(actual)
}

func (s *EngineSuite) Test113() {
	actual := NotNilFn(MessageTuple{}, NamedArgs{"Value": 0})
	s.True(actual)
}

func (s *EngineSuite) Test114() {
	actual := NotNilFn(MessageTuple{}, NamedArgs{"Value": 100500})
	s.True(actual)
}

func (s *EngineSuite) Test115() {
	mt := MessageTuple{
		Message{
			Payload: map[string]any{
				"action":  "my_action",
				"double":  3.33,
				"int":     1,
				"boolean": false,
				"voltage": nil,
			},
		},
	}
	s.True(NotNilFn(mt, NamedArgs{"Value": "$message.action"}))
	s.True(NotNilFn(mt, NamedArgs{"Value": "$message.double"}))
	s.True(NotNilFn(mt, NamedArgs{"Value": "$message.int"}))
	s.True(NotNilFn(mt, NamedArgs{"Value": "$message.boolean"}))
	s.False(NotNilFn(mt, NamedArgs{"Value": "$message.voltage"}))
	s.False(NotNilFn(mt, NamedArgs{"Value": "$message.nonexisting"}))
}

func (s *EngineSuite) Test120() {
	actual := ChangedFn(MessageTuple{}, NamedArgs{})
	s.False(actual)
}

func (s *EngineSuite) Test121() {
	actual := ChangedFn(
		MessageTuple{
			Message{DeviceId: "foo1"},
			Message{DeviceId: "foo2"},
		},
		NamedArgs{"Value": "$deviceId"},
	)
	s.True(actual)
}

func (s *EngineSuite) Test130() {
	m := Message{}

	v, _ := m.Get("foo")
	s.Nil(v)

	v, _ = m.Get("Payload")
	s.Nil(v)

	v, _ = m.Get("RawPayload")
	s.Nil(v)

	v, _ = m.Get("ChannelMeta")
	s.NotNil(v)
	s.IsType(ChannelMeta{}, v)

	v, _ = m.Get("ChannelType")
	s.Equal(CHANNEL_UNKNOWN, v)

	v, _ = m.Get("DeviceClass")
	s.Equal(DEVICE_CLASS_UNKNOWN, v)

	v, _ = m.Get("DeviceId")
	s.Equal(DeviceId(""), v)

	v, _ = m.Get("Timestamp")
	s.Equal(time.Time{}, v)
}

func (s *EngineSuite) Test131() {

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

	v, _ := m.Get("foo")
	s.Nil(v)

	v, _ = m.Get("Payload")
	s.Nil(v)

	v, _ = m.Get("RawPayload")
	s.Nil(v)

	v, _ = m.Get("ChannelMeta")
	s.NotNil(v)

	v, _ = m.Get("action")
	s.Equal("single_left", v)

	v, _ = m.Get("voltage")
	s.Equal(3.33, v)

	v, _ = m.Get("offline")
	s.Equal(false, v)

	v, _ = m.Get("ChannelType")
	s.Equal(CHANNEL_MQTT, v)

	v, _ = m.Get("DeviceClass")
	s.Equal(DEVICE_CLASS_ZIGBEE_DEVICE, v)

	v, _ = m.Get("DeviceId")
	s.Equal(DeviceId("0x00158d0004244bda"), v)

}

func (s *EngineSuite) Test140() {
	Start(&dummyservice{})
}

func (s *EngineSuite) Test141() {
	Stop()
}

func TestEngine(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}
