package conditions

import (
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type ConditionsSuite struct {
	suite.Suite
}

func (s *ConditionsSuite) SetupSuite() {
}

func (s *ConditionsSuite) TeardownSuite() {
}

func (s *ConditionsSuite) Test80() {
	actual, err := Equal(types.MessageTuple{}, types.Args{})
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test81() {
	actual, err := Equal(types.MessageTuple{}, types.Args{"Left": 1, "Right": 1})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test82() {
	actual, err := Equal(types.MessageTuple{}, types.Args{"Left": "one", "Right": "one"})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test83() {
	actual, err := Equal(
		types.MessageTuple{
			Curr: &types.Message{
				Payload: map[string]any{
					"action": "my_action",
				},
			},
		},
		types.Args{"Left": "$message.action", "Right": "my_action"},
	)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test84() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			DeviceId: "0x00158d0004244bda",
		},
	}
	args := types.Args{
		"Left":  "$deviceId",
		"Right": types.DeviceId("0x00158d0004244bda"),
	}
	actual, err := Equal(mt, args)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test90() {
	actual, err := NotEqual(types.MessageTuple{}, types.Args{})
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test91() {
	actual, err := NotEqual(types.MessageTuple{}, types.Args{"Left": 1, "Right": 1})
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test92() {
	actual, err := NotEqual(types.MessageTuple{}, types.Args{"Left": "one", "Right": "one"})
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test93() {
	mt := types.MessageTuple{
		Curr: &types.Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := types.Args{"Left": "$message.action", "Right": "my_action"}
	actual, err := NotEqual(mt, args)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test100() {
	actual, err := InList(types.MessageTuple{}, types.Args{})
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test101() {
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
	actual, err := InList(mt, args)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test102() {
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
	actual, err := InList(mt, args)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test103() {
	mt := types.MessageTuple{
		// types.Message{},
	}
	args := types.Args{
		"Value": "some1",
		"List":  []any{"some2", "some3"},
	}
	actual, err := InList(mt, args)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test105() {
	mt := types.MessageTuple{}
	args := types.Args{
		"List":  []any{types.DeviceId("0x00158d0004244bda")},
		"Value": types.DeviceId("0x00158d0004244bda"),
	}
	actual, err := InList(mt, args)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test110() {
	actual, err := NotNil(types.MessageTuple{}, types.Args{})
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test111() {
	actual, err := NotNil(types.MessageTuple{}, types.Args{"Value": "foo"})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test112() {
	actual, err := NotNil(types.MessageTuple{}, types.Args{"Value": false})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test113() {
	actual, err := NotNil(types.MessageTuple{}, types.Args{"Value": 0})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test114() {
	actual, err := NotNil(types.MessageTuple{}, types.Args{"Value": 100500})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test120() {
	actual, err := Changed(types.MessageTuple{}, types.Args{})
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test121() {
	actual, err := Changed(
		types.MessageTuple{
			Curr: &types.Message{DeviceId: "foo1"},
			Prev: &types.Message{DeviceId: "foo2"},
		},
		types.Args{"Value": "$deviceId"},
	)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test122() {
	actual, err := Changed(
		types.MessageTuple{
			Curr: &types.Message{DeviceId: "foo1"},
		},
		types.Args{"Value": "$deviceId"},
	)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test150() {
	defer func() { _ = recover() }()
	actual, err := ZigbeeDevice(types.MessageTuple{}, types.Args{})
	s.Nil(err)
	s.False(actual)
	s.Fail("expected to panic")
}

func (s *ConditionsSuite) Test151() {
	mt := types.MessageTuple{
		Curr: &types.Message{DeviceId: "0x00158d0004244bda", DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE},
	}
	actual, err := ZigbeeDevice(mt, types.Args{"List": []any{types.DeviceId("0x00158d0004244bda")}})
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test104() {
	// defer func() { _ = recover() }()
	args := types.Args{
		"Value": "some1",
		"List":  "some2",
	}
	res, err := InList(types.MessageTuple{}, args)
	s.NotNil(err)
	s.False(res)
	// s.Fail("expected to panic")
}

func (s *ConditionsSuite) Test115() {
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

	var res bool
	var err error

	res, err = NotNil(mt, types.Args{"Value": "$message.action"})
	s.True(res)
	s.Nil(err)
	res, err = NotNil(mt, types.Args{"Value": "$message.double"})
	s.True(res)
	s.Nil(err)
	res, err = NotNil(mt, types.Args{"Value": "$message.int"})
	s.True(res)
	s.Nil(err)
	res, err = NotNil(mt, types.Args{"Value": "$message.boolean"})
	s.True(res)
	s.Nil(err)
	res, err = NotNil(mt, types.Args{"Value": "$message.voltage"})
	s.False(res)
	s.Nil(err)
	res, err = NotNil(mt, types.Args{"Value": "$message.nonexisting"})
	s.False(res)
	s.NotNil(err)
}

func TestConditions(t *testing.T) {
	suite.Run(t, new(ConditionsSuite))
}
