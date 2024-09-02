package engine

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
	e engine
}

func (s *EngineSuite) SetupSuite() {
	s.e = NewEngine(NewOptions())
}

func (s *EngineSuite) Test10() {
	actual := s.e.matchesCondition(types.MessageTuple{}, Condition{}, Rule{}, "Test10")
	s.True(actual)
}

func (s *EngineSuite) Test11() {
	actual := s.e.matchesCondition(types.MessageTuple{}, Condition{
		Or: true,
		List: []Condition{
			{Fn: COND_EQUAL, Args: Args{"Left": true, "Right": false}},
			{Fn: COND_EQUAL, Args: Args{"Left": 1.11, "Right": 1.11}},
		},
	}, Rule{}, "Test11")
	s.True(actual)
}

func (s *EngineSuite) Test12() {
	actual := s.e.matchesCondition(types.MessageTuple{}, Condition{
		List: []Condition{
			{Fn: COND_EQUAL, Args: Args{"Left": true, "Right": false}},
			{Fn: COND_EQUAL, Args: Args{"Left": 1.11, "Right": 1.11}},
		},
	}, Rule{}, "Test12")
	s.False(actual)
}

func (s *EngineSuite) Test20() {
	// defer func() { _ = recover() }()
	s.False(s.e.invokeConditionFunc(types.MessageTuple{}, 0, nil, Rule{}, "Test20"))
	// s.Fail("expected to panic")
}

func (s *EngineSuite) Test30() {
	actual := s.e.matchesListSome(types.MessageTuple{}, []Condition{}, Rule{}, "Test30")
	s.False(actual)
}

func (s *EngineSuite) Test31() {
	actual := s.e.matchesListSome(types.MessageTuple{}, []Condition{
		{Fn: COND_EQUAL, Args: Args{"Left": 1, "Right": 1}},
		{Fn: COND_EQUAL, Args: Args{"Left": "foo", "Right": "bar"}},
	}, Rule{}, "Test31")
	s.True(actual)
}

func (s *EngineSuite) Test40() {
	actual := s.e.matchesListEvery(types.MessageTuple{}, []Condition{}, Rule{}, "Test40")
	s.False(actual)
}

func (s *EngineSuite) Test41() {
	actual := s.e.matchesListEvery(types.MessageTuple{}, []Condition{
		{Fn: COND_EQUAL, Args: Args{"Left": 1, "Right": 1}},
		{Fn: COND_EQUAL, Args: Args{"Left": "foo", "Right": "foo"}},
	}, Rule{}, "Test41")
	s.True(actual)
}

func (s *EngineSuite) Test50() {
	actual, err := arg_value(types.Message{}, Args{"foo": 1}, "foo")
	expected := 1
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test51() {
	actual, err := arg_value(types.Message{DeviceId: "foo2"}, Args{"Lorem": "$message.DeviceId"}, "Lorem")
	expected := types.DeviceId("foo2")
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test52() {
	actual, err := arg_value(types.Message{DeviceId: "foo3"}, Args{"Lorem1": "$deviceId"}, "Lorem1")
	expected := types.DeviceId("foo3")
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test53() {
	m := types.Message{
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
	m := types.Message{
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
	m := types.Message{
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
	actual, err := arg_value(types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}, Args{"Lorem1": "$deviceClass"}, "Lorem1")
	s.Equal(types.DEVICE_CLASS_ZIGBEE_BRIDGE, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test60() {
	s.e.executeActions([]types.Message{}, Rule{}, "Test60")
}

func (s *EngineSuite) Test70() {
	s.e.handleMessage(types.Message{}, []Rule{})
}

func (s *EngineSuite) Test71() {
	s.e.handleMessage(types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}, []Rule{})
}

func (s *EngineSuite) Test72() {
	s.e.handleMessage(types.Message{}, []Rule{
		{
			Condition: Condition{
				Fn:   COND_EQUAL,
				Args: Args{"Left": true, "Right": true},
			},
		},
	})
}

// func (s *EngineSuite) Test73() {
// 	defer func() { _ = recover() }()
// 	s.e.handleMessage(types.Message{}, []Rule{
// 		{
// 			Condition: Condition{
// 				Fn:   COND_EQUAL,
// 				Args: Args{"Left": true, "Right": true},
// 			},
// 			Throttle: time.Minute,
// 		},
// 	})
// 	s.Fail("expected to panic")
// }

func (s *EngineSuite) Test80() {
	actual := Equal(types.MessageTuple{}, Args{}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test81() {
	actual := Equal(types.MessageTuple{}, Args{"Left": 1, "Right": 1}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test82() {
	actual := Equal(types.MessageTuple{}, Args{"Left": "one", "Right": "one"}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test83() {
	actual := Equal(
		types.MessageTuple{
			types.Message{
				Payload: map[string]any{
					"action": "my_action",
				},
			},
		},
		Args{"Left": "$message.action", "Right": "my_action"},
		&s.e,
	)
	s.True(actual)
}

func (s *EngineSuite) Test84() {
	mt := types.MessageTuple{
		types.Message{
			DeviceId: "0x00158d0004244bda",
		},
	}
	args := Args{
		"Left":  "$deviceId",
		"Right": types.DeviceId("0x00158d0004244bda"),
	}
	actual := Equal(mt, args, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test90() {
	actual := NotEqual(types.MessageTuple{}, Args{}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test91() {
	actual := NotEqual(types.MessageTuple{}, Args{"Left": 1, "Right": 1}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test92() {
	actual := NotEqual(types.MessageTuple{}, Args{"Left": "one", "Right": "one"}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test93() {
	mt := types.MessageTuple{
		types.Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := Args{"Left": "$message.action", "Right": "my_action"}
	actual := NotEqual(
		mt,
		args,
		&s.e,
	)
	s.False(actual)
}

func (s *EngineSuite) Test100() {
	actual := InList(types.MessageTuple{}, Args{}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test101() {
	mt := types.MessageTuple{
		types.Message{
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
	actual := InList(mt, args, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test102() {
	mt := types.MessageTuple{
		types.Message{
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
	actual := InList(mt, args, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test103() {
	mt := types.MessageTuple{
		// types.Message{},
	}
	args := Args{
		"Value": "some1",
		"List":  []any{"some2", "some3"},
	}
	actual := InList(mt, args, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test104() {
	defer func() { _ = recover() }()
	args := Args{
		"Value": "some1",
		"List":  "some2",
	}
	InList(types.MessageTuple{}, args, &s.e)
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test105() {
	mt := types.MessageTuple{}
	args := Args{
		"List":  []any{types.DeviceId("0x00158d0004244bda")},
		"Value": types.DeviceId("0x00158d0004244bda"),
	}
	actual := InList(mt, args, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test110() {
	actual := NotNil(types.MessageTuple{}, Args{}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test111() {
	actual := NotNil(types.MessageTuple{}, Args{"Value": "foo"}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test112() {
	actual := NotNil(types.MessageTuple{}, Args{"Value": false}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test113() {
	actual := NotNil(types.MessageTuple{}, Args{"Value": 0}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test114() {
	actual := NotNil(types.MessageTuple{}, Args{"Value": 100500}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test115() {
	mt := types.MessageTuple{
		types.Message{
			Payload: map[string]any{
				"action":  "my_action",
				"double":  3.33,
				"int":     1,
				"boolean": false,
				"voltage": nil,
			},
		},
	}
	s.True(NotNil(mt, Args{"Value": "$message.action"}, &s.e))
	s.True(NotNil(mt, Args{"Value": "$message.double"}, &s.e))
	s.True(NotNil(mt, Args{"Value": "$message.int"}, &s.e))
	s.True(NotNil(mt, Args{"Value": "$message.boolean"}, &s.e))
	s.False(NotNil(mt, Args{"Value": "$message.voltage"}, &s.e))
	s.False(NotNil(mt, Args{"Value": "$message.nonexisting"}, &s.e))
}

func (s *EngineSuite) Test120() {
	actual := Changed(types.MessageTuple{}, Args{}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test121() {
	actual := Changed(
		types.MessageTuple{
			types.Message{DeviceId: "foo1"},
			types.Message{DeviceId: "foo2"},
		},
		Args{"Value": "$deviceId"},
		&s.e,
	)
	s.True(actual)
}

func (s *EngineSuite) Test130() {
	m := types.Message{}

	v, _ := m.Get("foo")
	s.Nil(v)

	v, _ = m.Get("Payload")
	s.Nil(v)

	v, _ = m.Get("RawPayload")
	s.Nil(v)

	v, _ = m.Get("ChannelType")
	s.Equal(types.CHANNEL_UNKNOWN, v)

	v, _ = m.Get("DeviceClass")
	s.Equal(types.DEVICE_CLASS_UNKNOWN, v)

	v, _ = m.Get("DeviceId")
	s.Equal(types.DeviceId(""), v)

	v, _ = m.Get("Timestamp")
	s.Equal(time.Time{}, v)
}

func (s *EngineSuite) Test131() {

	m := types.Message{
		ChannelMeta: types.ChannelMeta{MqttTopic: "foo/111"},
		ChannelType: types.CHANNEL_MQTT,
		DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE,
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

	v, _ = m.Get("action")
	s.Equal("single_left", v)

	v, _ = m.Get("voltage")
	s.Equal(3.33, v)

	v, _ = m.Get("offline")
	s.Equal(false, v)

	v, _ = m.Get("ChannelType")
	s.Equal(types.CHANNEL_MQTT, v)

	v, _ = m.Get("DeviceClass")
	s.Equal(types.DEVICE_CLASS_ZIGBEE_DEVICE, v)

	v, _ = m.Get("DeviceId")
	s.Equal(types.DeviceId("0x00158d0004244bda"), v)

}

func (s *EngineSuite) Test140() {
	s.e.Start()
}

func (s *EngineSuite) Test141() {
	s.e.Stop()
}

func (s *EngineSuite) Test142() {
	opts := NewOptions()
	opts.SetLogTag(func(m string) string { return " " })
}

func (s *EngineSuite) Test150() {
	actual := ZigbeeDevice(types.MessageTuple{}, Args{}, &s.e)
	s.False(actual)
}

func (s *EngineSuite) Test151() {
	mt := types.MessageTuple{types.Message{DeviceId: "0x00158d0004244bda", DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE}}
	actual := ZigbeeDevice(mt, Args{"List": []any{types.DeviceId("0x00158d0004244bda")}}, &s.e)
	s.True(actual)
}

func (s *EngineSuite) Test160() {
	input := []byte(`{"Foo":1}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
}

func (s *EngineSuite) Test161() {
	input := []byte(`{"Foo":"bar"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
	s.Equal(args["Foo"], "bar")
}

func (s *EngineSuite) Test162() {
	input := []byte(`{"Lorem":"DeviceId(bar-111)"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Lorem")
	s.Equal(args["Lorem"], types.DeviceId("bar-111"))
	s.IsType(args["Lorem"], types.DeviceId(""))
	fmt.Println(args)
}

func (s *EngineSuite) Test163() {
	input := []byte(`{"ClassesList":["DeviceClass(1)","DeviceClass(2)"]}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "ClassesList")
	s.Len(args["ClassesList"], 2)
	s.Equal("map[ClassesList:[zigbee-device (id=1) device-pinger (id=2)]]", fmt.Sprintf("%v", args))
}

func (s *EngineSuite) Test164() {
	input := []byte(`{foo}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	fmt.Println(args)
	s.NotNil(err)
}

func (s *EngineSuite) Test165() {
	input := []byte(`{"foo":"DeviceId"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Equal("DeviceId", args["foo"])
}

func (s *EngineSuite) Test170() {
	args := Args{"Foo1": types.DEVICE_CLASS_BOT}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo1":5}`, string(argsjson))
}

func (s *EngineSuite) Test171() {
	args := Args{"Foo2": types.DeviceId("some-111")}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo2":"some-111"}`, string(argsjson))
}

func TestEngine(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}
