package engine

import (
	"encoding/json"
	"fmt"
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	conditions "github.com/fedulovivan/mhz19-go/internal/engine/conditions"
	"github.com/fedulovivan/mhz19-go/internal/entities/ldm"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type EngineSuite struct {
	suite.Suite
	e types.Engine
}

func (s *EngineSuite) SetupSuite() {
	s.e = NewEngine()
	s.e.SetLdmService(ldm.NewService(ldm.RepoSingleton()))
}

func (s *EngineSuite) Test10() {
	actual := s.e.MatchesCondition(types.MessageTuple{}, types.Condition{}, types.Rule{}, "Test10")
	s.True(actual)
}

func (s *EngineSuite) Test11() {
	actual := s.e.MatchesCondition(types.MessageTuple{}, types.Condition{
		Or: true,
		List: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": false}},
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1.11, "Right": 1.11}},
		},
	}, types.Rule{}, "Test11")
	s.True(actual)
}

func (s *EngineSuite) Test12() {
	actual := s.e.MatchesCondition(types.MessageTuple{}, types.Condition{
		List: []types.Condition{
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": true, "Right": false}},
			{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1.11, "Right": 1.11}},
		},
	}, types.Rule{}, "Test12")
	s.False(actual)
}

func (s *EngineSuite) Test20() {
	defer func() { _ = recover() }()
	s.False(s.e.InvokeConditionFunc(types.MessageTuple{}, 0, nil, types.Rule{}, "Test20"))
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test30() {
	actual := s.e.MatchesListSome(types.MessageTuple{}, []types.Condition{}, types.Rule{}, "Test30")
	s.False(actual)
}

func (s *EngineSuite) Test31() {
	actual := s.e.MatchesListSome(types.MessageTuple{}, []types.Condition{
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": "foo", "Right": "bar"}},
	}, types.Rule{}, "Test31")
	s.True(actual)
}

func (s *EngineSuite) Test40() {
	actual := s.e.MatchesListEvery(types.MessageTuple{}, []types.Condition{}, types.Rule{}, "Test40")
	s.False(actual)
}

func (s *EngineSuite) Test41() {
	actual := s.e.MatchesListEvery(types.MessageTuple{}, []types.Condition{
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": 1, "Right": 1}},
		{Fn: types.COND_EQUAL, Args: types.Args{"Left": "foo", "Right": "foo"}},
	}, types.Rule{}, "Test41")
	s.True(actual)
}

func (s *EngineSuite) Test50() {
	r := arg_reader.NewArgReader(nil, types.Args{"foo": 1}, nil, nil, nil)
	actual, err := r.Stage1("foo")
	expected := 1
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test51() {
	r := arg_reader.NewArgReader(&types.Message{DeviceId: "foo2"}, types.Args{"Lorem": "$message.DeviceId"}, nil, nil, nil)
	actual, err := r.Stage1("Lorem")
	expected := types.DeviceId("foo2")
	s.Equal(expected, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test52() {
	r := arg_reader.NewArgReader(&types.Message{DeviceId: "foo3"}, types.Args{"Lorem1": "$deviceId"}, nil, nil, nil)
	actual, err := r.Stage1("Lorem1")
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
	args := types.Args{
		"Lorem3": "$message.action",
	}
	r := arg_reader.NewArgReader(&m, args, nil, nil, nil)
	actual, err := r.Stage1("Lorem3")
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
	args := types.Args{
		"Lorem3": "$message.voltage",
	}
	r := arg_reader.NewArgReader(&m, args, nil, nil, nil)
	actual, err := r.Stage1("Lorem3")
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
	args := types.Args{
		"Lorem3": "$message.bar",
	}
	r := arg_reader.NewArgReader(&m, args, nil, nil, nil)
	actual, err := r.Stage1("Lorem3")
	s.Nil(actual)
	s.NotEmpty(err)
}

func (s *EngineSuite) Test56() {
	m := types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}
	args := types.Args{"Lorem1": "$deviceClass"}
	r := arg_reader.NewArgReader(&m, args, nil, nil, nil)
	actual, err := r.Stage1("Lorem1")
	s.Equal(types.DEVICE_CLASS_ZIGBEE_BRIDGE, actual)
	s.Nil(err)
}

func (s *EngineSuite) Test60() {
	s.e.ExecuteActions([]types.Message{}, types.Rule{}, "Test60")
}

func (s *EngineSuite) Test70() {
	defer func() { _ = recover() }()
	s.e.HandleMessage(types.Message{}, []types.Rule{})
	s.e.HandleMessage(types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}, []types.Rule{})
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test72() {
	s.e.HandleMessage(types.Message{}, []types.Rule{
		{
			Condition: types.Condition{
				Fn:   types.COND_EQUAL,
				Args: types.Args{"Left": true, "Right": true},
			},
		},
	})
}

func (s *EngineSuite) Test80() {
	actual := conditions.Equal(types.MessageTuple{}, types.Args{}, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test81() {
	actual := conditions.Equal(types.MessageTuple{}, types.Args{"Left": 1, "Right": 1}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test82() {
	actual := conditions.Equal(types.MessageTuple{}, types.Args{"Left": "one", "Right": "one"}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test83() {
	actual := conditions.Equal(
		types.MessageTuple{
			Curr: &types.Message{
				Payload: map[string]any{
					"action": "my_action",
				},
			},
		},
		types.Args{"Left": "$message.action", "Right": "my_action"},
		s.e,
	)
	s.True(actual)
}

func (s *EngineSuite) Test84() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			DeviceId: "0x00158d0004244bda",
		},
	}
	args := types.Args{
		"Left":  "$deviceId",
		"Right": types.DeviceId("0x00158d0004244bda"),
	}
	actual := conditions.Equal(mt, args, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test90() {
	actual := conditions.NotEqual(types.MessageTuple{}, types.Args{}, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test91() {
	actual := conditions.NotEqual(types.MessageTuple{}, types.Args{"Left": 1, "Right": 1}, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test92() {
	actual := conditions.NotEqual(types.MessageTuple{}, types.Args{"Left": "one", "Right": "one"}, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test93() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := types.Args{"Left": "$message.action", "Right": "my_action"}
	actual := conditions.NotEqual(
		mt,
		args,
		s.e,
	)
	s.False(actual)
}

func (s *EngineSuite) Test100() {
	actual := conditions.InList(types.MessageTuple{}, types.Args{}, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test101() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := types.Args{
		"Value": "$message.action",
		"List": []any{
			"foo1",
			"my_action",
		},
	}
	actual := conditions.InList(mt, args, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test102() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			Payload: map[string]any{
				"voltage": 1.11,
			},
		},
	}
	args := types.Args{
		"Value": "$message.voltage",
		"List": []any{
			1,
			1.11,
			2.0,
		},
	}
	actual := conditions.InList(mt, args, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test103() {
	mt := types.MessageTuple{
		// types.Message{},
	}
	args := types.Args{
		"Value": "some1",
		"List":  []any{"some2", "some3"},
	}
	actual := conditions.InList(mt, args, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test104() {
	defer func() { _ = recover() }()
	args := types.Args{
		"Value": "some1",
		"List":  "some2",
	}
	conditions.InList(types.MessageTuple{}, args, s.e)
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test105() {
	mt := types.MessageTuple{}
	args := types.Args{
		"List":  []any{types.DeviceId("0x00158d0004244bda")},
		"Value": types.DeviceId("0x00158d0004244bda"),
	}
	actual := conditions.InList(mt, args, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test110() {
	actual := conditions.NotNil(types.MessageTuple{}, types.Args{}, s.e)
	s.False(actual)
}

func (s *EngineSuite) Test111() {
	actual := conditions.NotNil(types.MessageTuple{}, types.Args{"Value": "foo"}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test112() {
	actual := conditions.NotNil(types.MessageTuple{}, types.Args{"Value": false}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test113() {
	actual := conditions.NotNil(types.MessageTuple{}, types.Args{"Value": 0}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test114() {
	actual := conditions.NotNil(types.MessageTuple{}, types.Args{"Value": 100500}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test115() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			Payload: map[string]any{
				"action":  "my_action",
				"double":  3.33,
				"int":     1,
				"boolean": false,
				"voltage": nil,
			},
		},
	}
	s.True(conditions.NotNil(mt, types.Args{"Value": "$message.action"}, s.e))
	s.True(conditions.NotNil(mt, types.Args{"Value": "$message.double"}, s.e))
	s.True(conditions.NotNil(mt, types.Args{"Value": "$message.int"}, s.e))
	s.True(conditions.NotNil(mt, types.Args{"Value": "$message.boolean"}, s.e))
	s.False(conditions.NotNil(mt, types.Args{"Value": "$message.voltage"}, s.e))
	s.False(conditions.NotNil(mt, types.Args{"Value": "$message.nonexisting"}, s.e))
}

func (s *EngineSuite) Test120() {
	// defer func() { _ = recover() }()
	actual := conditions.Changed(types.MessageTuple{}, types.Args{}, s.e)
	s.False(actual)
	// s.Fail("expected to panic")
}

func (s *EngineSuite) Test121() {
	actual := conditions.Changed(
		types.MessageTuple{
			Curr: &types.Message{DeviceId: "foo1"},
			Prev: &types.Message{DeviceId: "foo2"},
		},
		types.Args{"Value": "$deviceId"},
		s.e,
	)
	s.True(actual)
}

func (s *EngineSuite) Test122() {
	actual := conditions.Changed(
		types.MessageTuple{
			Curr: &types.Message{DeviceId: "foo1"},
		},
		types.Args{"Value": "$deviceId"},
		s.e,
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

func (s *EngineSuite) Test150() {
	defer func() { _ = recover() }()
	actual := conditions.ZigbeeDevice(types.MessageTuple{}, types.Args{}, s.e)
	s.False(actual)
	s.Fail("expected to panic")
}

func (s *EngineSuite) Test151() {
	mt := types.MessageTuple{
		Curr: &types.Message{DeviceId: "0x00158d0004244bda", DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE},
	}
	actual := conditions.ZigbeeDevice(mt, types.Args{"List": []any{types.DeviceId("0x00158d0004244bda")}}, s.e)
	s.True(actual)
}

func (s *EngineSuite) Test160() {
	input := []byte(`{"Foo":1}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
}

func (s *EngineSuite) Test161() {
	input := []byte(`{"Foo":"bar"}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
	s.Equal(args["Foo"], "bar")
}

func (s *EngineSuite) Test162() {
	input := []byte(`{"Lorem":"DeviceId(bar-111)"}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Lorem")
	s.Equal(args["Lorem"], types.DeviceId("bar-111"))
	s.IsType(args["Lorem"], types.DeviceId(""))
	fmt.Println(args)
}

func (s *EngineSuite) Test163() {
	input := []byte(`{"ClassesList":["DeviceClass(1)","DeviceClass(2)"]}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "ClassesList")
	s.Len(args["ClassesList"], 2)
	s.Equal("map[ClassesList:[zigbee-device (id=1) device-pinger (id=2)]]", fmt.Sprintf("%v", args))
}

func (s *EngineSuite) Test164() {
	input := []byte(`{foo}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	fmt.Println(args)
	s.NotNil(err)
}

func (s *EngineSuite) Test165() {
	input := []byte(`{"foo":"DeviceId"}`)
	args := types.Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Equal("DeviceId", args["foo"])
}

func (s *EngineSuite) Test170() {
	args := types.Args{"Foo1": types.DEVICE_CLASS_BOT}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo1":5}`, string(argsjson))
}

func (s *EngineSuite) Test171() {
	args := types.Args{"Foo2": types.DeviceId("some-111")}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo2":"some-111"}`, string(argsjson))
}

func TestEngine(t *testing.T) {
	suite.Run(t, new(EngineSuite))
}
