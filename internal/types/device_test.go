package types

import (
	"encoding/json"
	"testing"
	"time"

	"github.com/stretchr/testify/suite"
)

type DeviceSuite struct {
	suite.Suite
}

func (s *DeviceSuite) SetupSuite() {

}

func (s *DeviceSuite) TeardownSuite() {

}

func (s *DeviceSuite) Test10() {
	type withdeviceid struct {
		Field DeviceId
	}
	target := withdeviceid{}
	data := []byte(`{"Field":"DeviceId(foo)"}`)
	err := json.Unmarshal(data, &target)
	s.Nil(err)
	s.IsType(DeviceId("bar"), target.Field)
	s.Equal(DeviceId("foo"), target.Field)
}

func (s *DeviceSuite) Test11() {
	type withdeviceid struct {
		Field DeviceId
	}
	target := withdeviceid{}
	data := []byte(`{}`)
	err := json.Unmarshal(data, &target)
	s.Nil(err)
	s.Equal(DeviceId(""), target.Field)
}

func (s *DeviceSuite) Test12() {
	type withdeviceid struct {
		Field DeviceId
	}
	target := withdeviceid{}
	data := []byte(`{"Field":null}`)
	err := json.Unmarshal(data, &target)
	s.Nil(err)
	s.Equal(DeviceId(""), target.Field)
}

func (s *DeviceSuite) Test20() {
	var target DeviceId
	data := []byte(`"DeviceId(foo)"`)
	err := json.Unmarshal(data, &target)
	s.Nil(err)
	s.IsType(DeviceId("bar"), target)
	s.Equal(DeviceId("foo"), target)
}

func (s *DeviceSuite) Test30() {
	var target DeviceId
	data := []byte(`123`)
	err := json.Unmarshal(data, &target)
	s.ErrorContains(err, "json: cannot unmarshal number into Go value of type string")
}

func (s *DeviceSuite) Test40() {
	var target DeviceId
	data := []byte(`"666"`)
	err := json.Unmarshal(data, &target)
	s.ErrorContains(err, `cannot parse string "666" into DeviceId`)
}

func (s *DeviceSuite) Test50() {
	name := "lorem"
	comments := "iprum dolorovi4"
	origin := "unit-test"
	device := Device{
		Id:          13,
		DeviceId:    DeviceId("foo"),
		DeviceClass: DEVICE_CLASS_BOT,
		Name:        &name,
		Comments:    &comments,
		Origin:      &origin,
		Json: map[string]any{
			"foo": 666,
		},
		BuriedTimeout: &BuriedTimeout{time.Second},
	}
	bytes, err := json.Marshal(device)
	s.Nil(err)
	expected := `{"id":13,"deviceId":"DeviceId(foo)","deviceClass":"telegram-bot","name":"lorem","comments":"iprum dolorovi4","origin":"unit-test","json":{"foo":666},"buriedTimeout":"1s"}`
	s.Equal(expected, string(bytes))
}

func TestDevice(t *testing.T) {
	suite.Run(t, new(DeviceSuite))
}
