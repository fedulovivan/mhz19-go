package conditions

import (
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	"github.com/stretchr/testify/suite"
)

type ConditionsSuite struct {
	suite.Suite
	tag utils.Tag
}

func (s *ConditionsSuite) SetupSuite() {
}

func (s *ConditionsSuite) TeardownSuite() {
}

func (s *ConditionsSuite) Test80() {
	actual, err := Equal(types.MessageCompound{}, types.Args{}, s.tag)
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test81() {
	actual, err := Equal(types.MessageCompound{}, types.Args{"Left": 1, "Right": 1}, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test82() {
	actual, err := Equal(types.MessageCompound{}, types.Args{"Left": "one", "Right": "one"}, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test83() {
	actual, err := Equal(
		types.MessageCompound{
			Curr: &types.Message{
				Payload: map[string]any{
					"action": "my_action",
				},
			},
		},
		types.Args{"Left": "$message.action", "Right": "my_action"},
		s.tag,
	)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test84() {
	mt := types.MessageCompound{
		Curr: &types.Message{
			DeviceId: "0x00158d0004244bda",
		},
	}
	args := types.Args{
		"Left":  "$deviceId",
		"Right": types.DeviceId("0x00158d0004244bda"),
	}
	actual, err := Equal(mt, args, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test93() {
	mt := types.MessageCompound{
		Curr: &types.Message{
			Payload: map[string]any{
				"action": "my_action",
			},
		},
	}
	args := types.Args{"Left": "$message.action", "Right": "my_action"}
	actual, err := Equal(mt, args, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test100() {
	actual, err := InList(types.MessageCompound{}, types.Args{}, s.tag)
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test101() {
	mt := types.MessageCompound{
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
	actual, err := InList(mt, args, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test102() {
	mt := types.MessageCompound{
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
	actual, err := InList(mt, args, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test103() {
	mt := types.MessageCompound{
		// types.Message{},
	}
	args := types.Args{
		"Value": "some1",
		"List":  []any{"some2", "some3"},
	}
	actual, err := InList(mt, args, s.tag)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test105() {
	mt := types.MessageCompound{}
	args := types.Args{
		"List":  []any{types.DeviceId("0x00158d0004244bda")},
		"Value": types.DeviceId("0x00158d0004244bda"),
	}
	actual, err := InList(mt, args, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test110() {
	actual, err := Nil(types.MessageCompound{}, types.Args{}, s.tag)
	s.NotNil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test111() {
	actual, err := Nil(types.MessageCompound{}, types.Args{"Value": "foo"}, s.tag)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test112() {
	actual, err := Nil(types.MessageCompound{}, types.Args{"Value": false}, s.tag)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test113() {
	actual, err := Nil(types.MessageCompound{}, types.Args{"Value": 0}, s.tag)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test114() {
	actual, err := Nil(types.MessageCompound{}, types.Args{"Value": 100500}, s.tag)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test120() {
	actual, err := Changed(types.MessageCompound{}, types.Args{}, s.tag)
	s.NotNil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test121() {
	actual, err := Changed(
		types.MessageCompound{
			Curr: &types.Message{DeviceId: "foo1"},
			Prev: &types.Message{DeviceId: "foo2"},
		},
		types.Args{"Value": "$deviceId"},
		s.tag,
	)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test122() {
	actual, err := Changed(
		types.MessageCompound{
			Curr: &types.Message{DeviceId: "foo1"},
		},
		types.Args{"Value": "$deviceId"},
		s.tag,
	)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test150() {
	actual, err := ZigbeeDevice(types.MessageCompound{}, types.Args{}, s.tag)
	s.EqualError(err, "[]any is expected for List")
	s.False(actual)
}

func (s *ConditionsSuite) Test151() {
	mt := types.MessageCompound{
		Curr: &types.Message{DeviceId: "0x00158d0004244bda", DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE},
	}
	actual, err := ZigbeeDevice(mt, types.Args{"List": []any{types.DeviceId("0x00158d0004244bda")}}, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test104() {
	args := types.Args{
		"Value": "some1",
		"List":  "some2",
	}
	res, err := InList(types.MessageCompound{}, args, s.tag)
	s.NotNil(err)
	s.False(res)
}

func (s *ConditionsSuite) Test115() {
	mt := types.MessageCompound{
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
	res, err = Nil(mt, types.Args{"Value": "$message.action"}, s.tag)
	s.False(res)
	s.Nil(err)
	res, err = Nil(mt, types.Args{"Value": "$message.double"}, s.tag)
	s.False(res)
	s.Nil(err)
	res, err = Nil(mt, types.Args{"Value": "$message.int"}, s.tag)
	s.False(res)
	s.Nil(err)
	res, err = Nil(mt, types.Args{"Value": "$message.boolean"}, s.tag)
	s.False(res)
	s.Nil(err)
	res, err = Nil(mt, types.Args{"Value": "$message.voltage"}, s.tag)
	s.True(res)
	s.Nil(err)
	res, err = Nil(mt, types.Args{"Value": "$message.nonexisting"}, s.tag)
	s.True(res)
	s.Nil(err)
}

func (s *ConditionsSuite) Test160() {
	mt := types.MessageCompound{
		Curr: &types.Message{},
	}
	actual, err := Channel(mt, types.Args{"Value": types.ChannelType(0)}, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test161() {
	mt := types.MessageCompound{
		Curr: &types.Message{
			ChannelType: types.CHANNEL_TELEGRAM,
		},
	}
	actual, err := Channel(mt, types.Args{"Value": types.CHANNEL_TELEGRAM}, s.tag)
	s.Nil(err)
	s.True(actual)
}

func (s *ConditionsSuite) Test170() {
	mt := types.MessageCompound{
		Curr: &types.Message{},
	}
	actual, err := FromEndDevice(mt, types.Args{}, s.tag)
	s.Nil(err)
	s.False(actual)
}

func (s *ConditionsSuite) Test171() {
	mt := types.MessageCompound{
		Curr: &types.Message{
			FromEndDevice: true,
		},
	}
	actual, err := FromEndDevice(mt, types.Args{}, s.tag)
	s.Nil(err)
	s.True(actual)
}

func TestConditions(t *testing.T) {
	suite.Run(t, &ConditionsSuite{
		tag: utils.NewTag(logger.CONDS),
	})
}
