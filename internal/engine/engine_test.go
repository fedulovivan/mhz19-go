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
	actual := matchesCondition(MessageTuple{}, Condition{}, Rule{}, "Test10")
	s.False(actual)
}

func (s *EngineSuite) Test11() {
	actual := matchesCondition(MessageTuple{}, Condition{
		Or: true,
		List: []Condition{
			{Fn: COND_EQUAL, Args: Args{"Left": true, "Right": false}},
			{Fn: COND_EQUAL, Args: Args{"Left": 1.11, "Right": 1.11}},
		},
	}, Rule{}, "Test11")
	s.True(actual)
}

func (s *EngineSuite) Test12() {
	actual := matchesCondition(MessageTuple{}, Condition{
		List: []Condition{
			{Fn: COND_EQUAL, Args: Args{"Left": true, "Right": false}},
			{Fn: COND_EQUAL, Args: Args{"Left": 1.11, "Right": 1.11}},
		},
	}, Rule{}, "Test12")
	s.False(actual)
}

func (s *EngineSuite) Test20() {
	// defer func() { _ = recover() }()
	s.False(invokeConditionFunc(MessageTuple{}, 0, nil, Rule{}, "Test20"))
	// s.Fail("expected to panic")
}

func (s *EngineSuite) Test30() {
	actual := matchesListSome(MessageTuple{}, []Condition{}, Rule{}, "Test30")
	s.False(actual)
}

func (s *EngineSuite) Test31() {
	actual := matchesListSome(MessageTuple{}, []Condition{
		{Fn: COND_EQUAL, Args: Args{"Left": 1, "Right": 1}},
		{Fn: COND_EQUAL, Args: Args{"Left": "foo", "Right": "bar"}},
	}, Rule{}, "Test31")
	s.True(actual)
}

func (s *EngineSuite) Test40() {
	actual := matchesListEvery(MessageTuple{}, []Condition{}, Rule{}, "Test40")
	s.False(actual)
}

func (s *EngineSuite) Test41() {
	actual := matchesListEvery(MessageTuple{}, []Condition{
		{Fn: COND_EQUAL, Args: Args{"Left": 1, "Right": 1}},
		{Fn: COND_EQUAL, Args: Args{"Left": "foo", "Right": "foo"}},
	}, Rule{}, "Test41")
	s.True(actual)
}

func (s *EngineSuite) Test50() {
	actual, err := arg_value(Message{}, Args{"foo": 1}, "foo")
	expected := 1
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test51() {
	actual, err := arg_value(Message{DeviceId: "foo2"}, Args{"Lorem": "$message.DeviceId"}, "Lorem")
	expected := DeviceId("foo2")
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test52() {
	actual, err := arg_value(Message{DeviceId: "foo3"}, Args{"Lorem1": "$deviceId"}, "Lorem1")
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
	args := Args{
		"Lorem3": "$message.action",
	}
	actual, err := arg_value(m, args, "Lorem3")
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
	args := Args{
		"Lorem3": "$message.voltage",
	}
	actual, err := arg_value(m, args, "Lorem3")
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
	args := Args{
		"Lorem3": "$message.bar",
	}
	actual, err := arg_value(m, args, "Lorem3")
	s.Nil(actual)
	s.NotEmpty(err)
}

func (s *EngineSuite) Test56() {
	actual, err := arg_value(Message{DeviceClass: DEVICE_CLASS_ZIGBEE_BRIDGE}, Args{"Lorem1": "$deviceClass"}, "Lorem1")
	s.Equal(DEVICE_CLASS_ZIGBEE_BRIDGE, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test60() {
	executeActions([]Message{}, []Action{}, Rule{}, "Test60")
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
				Fn:   COND_EQUAL,
				Args: Args{"Left": true, "Right": true},
			},
		},
	})
}

func (s *EngineSuite) Test73() {
	defer func() { _ = recover() }()
	handleMessage(Message{}, []Rule{
		{
			Condition: Condition{
				Fn:   COND_EQUAL,
				Args: Args{"Left": true, "Right": true},
			},
			Throttle: time.Minute,
		},
	})
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test80() {
	actual := Equal(MessageTuple{}, Args{})
	s.False(actual)
}

func (s *EngineSuite) Test81() {
	actual := Equal(MessageTuple{}, Args{"Left": 1, "Right": 1})
	s.True(actual)
}

func (s *EngineSuite) Test82() {
	actual := Equal(MessageTuple{}, Args{"Left": "one", "Right": "one"})
	s.True(actual)
}

func (s *EngineSuite) Test83() {
	actual := Equal(
		MessageTuple{
			Message{
				Payload: map[string]any{
					"action": "my_action",
				},
			},
		},
		Args{"Left": "$message.action", "Right": "my_action"},
	)
	s.True(actual)
}

func (s *EngineSuite) Test84() {
	mt := MessageTuple{
		Message{
			DeviceId: "0x00158d0004244bda",
		},
	}
	args := Args{
		"Left":  "$deviceId",
		"Right": DeviceId("0x00158d0004244bda"),
	}
	actual := Equal(mt, args)
	s.True(actual)
}

func (s *EngineSuite) Test90() {
	actual := NotEqual(MessageTuple{}, Args{})
	s.False(actual)
}

func (s *EngineSuite) Test91() {
	actual := NotEqual(MessageTuple{}, Args{"Left": 1, "Right": 1})
	s.False(actual)
}

func (s *EngineSuite) Test92() {
	actual := NotEqual(MessageTuple{}, Args{"Left": "one", "Right": "one"})
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
	args := Args{"Left": "$message.action", "Right": "my_action"}
	actual := NotEqual(
		mt,
		args,
	)
	s.False(actual)
}

func (s *EngineSuite) Test100() {
	actual := InList(MessageTuple{}, Args{})
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
	args := Args{
		"Value": "$message.action",
		"List": []any{
			"foo1",
			"my_action",
		},
	}
	actual := InList(mt, args)
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
	args := Args{
		"Value": "$message.voltage",
		"List": []any{
			1,
			1.11,
			2.0,
		},
	}
	actual := InList(mt, args)
	s.True(actual)
}

func (s *EngineSuite) Test103() {
	mt := MessageTuple{
		// Message{},
	}
	args := Args{
		"Value": "some1",
		"List":  []any{"some2", "some3"},
	}
	actual := InList(mt, args)
	s.False(actual)
}

func (s *EngineSuite) Test104() {
	defer func() { _ = recover() }()
	args := Args{
		"Value": "some1",
		"List":  "some2",
	}
	InList(MessageTuple{}, args)
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test105() {
	mt := MessageTuple{}
	args := Args{
		"List":  []any{DeviceId("0x00158d0004244bda")},
		"Value": DeviceId("0x00158d0004244bda"),
	}
	actual := InList(mt, args)
	s.True(actual)
}

func (s *EngineSuite) Test110() {
	actual := NotNil(MessageTuple{}, Args{})
	s.False(actual)
}

func (s *EngineSuite) Test111() {
	actual := NotNil(MessageTuple{}, Args{"Value": "foo"})
	s.True(actual)
}

func (s *EngineSuite) Test112() {
	actual := NotNil(MessageTuple{}, Args{"Value": false})
	s.True(actual)
}

func (s *EngineSuite) Test113() {
	actual := NotNil(MessageTuple{}, Args{"Value": 0})
	s.True(actual)
}

func (s *EngineSuite) Test114() {
	actual := NotNil(MessageTuple{}, Args{"Value": 100500})
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
	s.True(NotNil(mt, Args{"Value": "$message.action"}))
	s.True(NotNil(mt, Args{"Value": "$message.double"}))
	s.True(NotNil(mt, Args{"Value": "$message.int"}))
	s.True(NotNil(mt, Args{"Value": "$message.boolean"}))
	s.False(NotNil(mt, Args{"Value": "$message.voltage"}))
	s.False(NotNil(mt, Args{"Value": "$message.nonexisting"}))
}

func (s *EngineSuite) Test120() {
	actual := Changed(MessageTuple{}, Args{})
	s.False(actual)
}

func (s *EngineSuite) Test121() {
	actual := Changed(
		MessageTuple{
			Message{DeviceId: "foo1"},
			Message{DeviceId: "foo2"},
		},
		Args{"Value": "$deviceId"},
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

	// v, _ = m.Get("ChannelMeta")
	// s.NotNil(v)
	// s.IsType(ChannelMeta{}, v)

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

	// v, _ = m.Get("ChannelMeta")
	// s.NotNil(v)

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
	opts := NewOptions()
	opts.SetServices(&service{})
	Start(opts)
}

func (s *EngineSuite) Test141() {
	Stop()
}

func (s *EngineSuite) Test142() {
	opts := NewOptions()
	opts.SetLogTag(func(m string) string { return " " })
}

func (s *EngineSuite) Test150() {
	actual := ZigbeeDevice(MessageTuple{}, Args{})
	s.False(actual)
}

func (s *EngineSuite) Test151() {
	mt := MessageTuple{Message{DeviceId: "0x00158d0004244bda", DeviceClass: DEVICE_CLASS_ZIGBEE_DEVICE}}
	actual := ZigbeeDevice(mt, Args{"List": []any{DeviceId("0x00158d0004244bda")}})
	s.True(actual)
}

func TestEngine(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}
