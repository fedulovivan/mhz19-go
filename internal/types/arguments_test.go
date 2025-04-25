package types

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/stretchr/testify/suite"
)

type ArgumentsSuite struct {
	suite.Suite
}

func (s *ArgumentsSuite) SetupSuite() {

}

func (s *ArgumentsSuite) TeardownSuite() {

}

func (s *ArgumentsSuite) Test10() {
	input := []byte(`{"Foo":1}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
}

func (s *ArgumentsSuite) Test15() {
	input := []byte(`{`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.NotNil(err)
}

func (s *ArgumentsSuite) Test20() {
	input := []byte(`{"Foo":"bar"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Foo")
	s.Equal(args["Foo"], "bar")
}

func (s *ArgumentsSuite) Test30() {
	input := []byte(`{"Lorem":"DeviceId(bar-111)"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Lorem")
	s.Equal(args["Lorem"], DeviceId("bar-111"))
	s.IsType(args["Lorem"], DeviceId(""))
	fmt.Println(args)
}

func (s *ArgumentsSuite) Test40() {
	input := []byte(`{"ClassesList":["DeviceClass(1)","DeviceClass(2)"]}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "ClassesList")
	s.Len(args["ClassesList"], 2)
	s.Equal("map[ClassesList:[zigbee-device device-pinger]]", fmt.Sprintf("%v", args))
}

func (s *ArgumentsSuite) Test41() {
	input := []byte(`{"Value": "DeviceClass(zigbee-device)"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Value")
	s.Equal("map[Value:zigbee-device]", fmt.Sprintf("%v", args))
}

func (s *ArgumentsSuite) Test44() {
	input := []byte(`{"ClassesList":["DeviceClass(foo)"]}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.NotNil(err)
}

func (s *ArgumentsSuite) Test45() {
	input := []byte(`{"Value": "ChannelType(4)"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Value")
	s.Equal("map[Value:system]", fmt.Sprintf("%v", args))
}

func (s *ArgumentsSuite) Test46() {
	input := []byte(`{"Value": "ChannelType(system)"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Contains(args, "Value")
	s.Equal("map[Value:system]", fmt.Sprintf("%v", args))
}

func (s *ArgumentsSuite) Test47() {
	input := []byte(`{"Value": "ChannelType(qwe)"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.NotNil(err)
}

func (s *ArgumentsSuite) Test50() {
	input := []byte(`{foo}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	fmt.Println(args)
	s.NotNil(err)
}

func (s *ArgumentsSuite) Test60() {
	input := []byte(`{"foo":"DeviceId"}`)
	args := Args{}
	err := json.Unmarshal(input, &args)
	s.Nil(err)
	s.Equal("DeviceId", args["foo"])
}

func (s *ArgumentsSuite) Test70() {
	args := Args{"Foo1": DEVICE_CLASS_BOT}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo1":"telegram-bot"}`, string(argsjson))
}

func (s *ArgumentsSuite) Test80() {
	args := Args{"Foo2": DeviceId("some-111")}
	argsjson, err := json.Marshal(args)
	s.Nil(err)
	s.Equal(`{"Foo2":"DeviceId(some-111)"}`, string(argsjson))
}

func TestArguments(t *testing.T) {
	suite.Run(t, new(ArgumentsSuite))
}
