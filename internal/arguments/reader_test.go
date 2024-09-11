package arguments

import (
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/mocks"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type ReaderSuite struct {
	suite.Suite
}

func (s *ReaderSuite) SetupSuite() {
}

func (s *ReaderSuite) TeardownSuite() {
}

func (s *ReaderSuite) Test10() {
	r := NewReader(nil, types.Args{"foo": 1}, nil, nil, nil)
	actual := r.Get("foo")
	expected := 1
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test30() {
	r := NewReader(&types.Message{DeviceId: "foo3"}, types.Args{"Lorem1": "$deviceId"}, nil, nil, nil)
	actual := r.Get("Lorem1")
	expected := types.DeviceId("foo3")
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test40() {
	m := types.Message{
		Payload: map[string]any{
			"action": "my_action",
		},
	}
	args := types.Args{
		"Lorem3": "$message.action",
	}
	r := NewReader(&m, args, nil, nil, nil)
	actual := r.Get("Lorem3")
	expected := "my_action"
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test50() {
	m := types.Message{
		Payload: map[string]any{
			"voltage": 3.33,
		},
	}
	args := types.Args{
		"Lorem3": "$message.voltage",
	}
	r := NewReader(&m, args, nil, nil, nil)
	actual := r.Get("Lorem3")
	expected := 3.33
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test60() {
	m := types.Message{
		Payload: map[string]any{
			"foo": false,
		},
	}
	args := types.Args{
		"Lorem3": "$message.bar",
	}
	r := NewReader(&m, args, nil, nil, nil)
	actual := r.Get("Lorem3")
	s.Equal("$message.bar", actual)
	s.NotEmpty(r.Error())
}

func (s *ReaderSuite) Test70() {
	m := types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}
	args := types.Args{"Lorem1": "$deviceClass"}
	r := NewReader(&m, args, nil, nil, nil)
	actual := r.Get("Lorem1")
	s.Equal(types.DEVICE_CLASS_ZIGBEE_BRIDGE, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test80() {
	m := types.Message{}
	args := types.Args{}
	mapping := types.Mapping{}
	r := NewReader(&m, args, mapping, nil, nil)
	s.NotNil(r)
}

func (s *ReaderSuite) Test90() {
	args := types.Args{}
	r := NewReader(nil, args, nil, nil, nil)
	s.NotNil(r)
	r.Get("Foo")
	s.EqualError(r.Error(), "no such argument Foo")
	r.Get("$message.bar")
	s.EqualError(r.Error(), "no such argument Foo\nno such argument $message.bar")
}

func (s *ReaderSuite) Test100() {
	m := types.Message{}
	args := types.Args{
		"Foo1": "{{ range .Messages }}{{end}}",
	}
	r := NewReader(&m, args, nil, &types.TemplatePayload{Messages: []types.Message{}}, nil)
	s.Equal("", r.Get("Foo1"))
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test101() {
	m := types.Message{}
	args := types.Args{
		"Foo2": "{{ .Test }}",
	}
	r := NewReader(&m, args, nil, &types.TemplatePayload{Messages: []types.Message{}}, nil)
	s.Equal(r.Get("Foo2"), "{{ .Test }}")
	s.EqualError(r.Error(), `template: Foo2:1:3: executing "Foo2" at <.Test>: can't evaluate field Test in type *types.TemplatePayload`)
}

func (s *ReaderSuite) Test102() {
	m := types.Message{}
	args := types.Args{
		"Foo2": "{{ }} }}",
	}
	r := NewReader(&m, args, nil, &types.TemplatePayload{Messages: []types.Message{}}, nil)
	s.Equal(r.Get("Foo2"), "{{ }} }}")
	s.EqualError(r.Error(), `template: Foo2:1: missing value for command`)
}

func (s *ReaderSuite) Test110() {
	m := types.Message{}
	args := types.Args{"Foo3": "lorem"}
	r := NewReader(&m, args, types.Mapping{"Foo3": {"lorem": "dolor"}}, nil, nil)
	s.Equal("dolor", r.Get("Foo3"))
}

func (s *ReaderSuite) Test111() {
	m := types.Message{}
	args := types.Args{"Foo4": 1}
	r := NewReader(&m, args, types.Mapping{"Foo4": {"1": "2"}}, nil, nil)
	s.Equal("2", r.Get("Foo4"))
}

func (s *ReaderSuite) Test120() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ range .Messages }}{{ .DeviceId }} - {{ deviceName .DeviceId }}{{end}}"}
	tpayload := types.TemplatePayload{Messages: []types.Message{{
		DeviceId: "10011cec96",
	}}}
	r := NewReader(&m, args, nil, &tpayload, engine)
	s.Equal("10011cec96 - My perfect name", r.Get("Foo5"))
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test121() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ range .Messages }}{{ .DeviceId }} - {{ deviceName .DeviceId }}{{end}}"}
	tpayload := types.TemplatePayload{Messages: []types.Message{{
		DeviceId: "lorem111",
	}}}
	r := NewReader(&m, args, nil, &tpayload, engine)
	s.Equal("lorem111 - lorem111", r.Get("Foo5"))
	// s.EqualError(r.Error(), `template: Foo5:1:42: executing "Foo5" at <deviceName .DeviceId>: error calling deviceName: no such device`)
}

func (s *ReaderSuite) Test122() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ range .Messages }}device name is {{ deviceName .DeviceId }}{{end}}"}
	tpayload := types.TemplatePayload{Messages: []types.Message{{
		DeviceId: "nullish-device-id",
	}}}
	r := NewReader(&m, args, nil, &tpayload, engine)
	expected := "device name is <unknonwn device originated from > nullish-device-id"
	actual := r.Get("Foo5")
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test123() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ deviceName 111 }}"}
	tpayload := types.TemplatePayload{}
	r := NewReader(&m, args, nil, &tpayload, engine)
	r.Get("Foo5")
	s.EqualError(r.Error(), `template: Foo5:1:3: executing "Foo5" at <deviceName 111>: error calling deviceName: deviceName accepts only types.DeviceId as an argument`)
}

func (s *ReaderSuite) Test124() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{
		"Foo5": "{{ pingerStatusName 1 }}",
		"Foo6": "{{ pingerStatusName 0 }}",
	}
	tpayload := types.TemplatePayload{}
	r := NewReader(&m, args, nil, &tpayload, engine)
	expected1 := "ONLINE"
	actual1 := r.Get("Foo5")
	s.Equal(expected1, actual1)
	expected2 := "OFFLINE"
	actual2 := r.Get("Foo6")
	s.Equal(expected2, actual2)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test125() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{
		"Foo5": "{{ pingerStatusName 333 }}",
	}
	tpayload := types.TemplatePayload{}
	r := NewReader(&m, args, nil, &tpayload, engine)
	expected := "UNKNOWN"
	actual := r.Get("Foo5")
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func TestReader(t *testing.T) {
	suite.Run(t, new(ReaderSuite))
}

// func (s *ReaderSuite) Test20() {
// 	r := NewReader(&types.Message{DeviceId: "foo2"}, types.Args{"Lorem": "$message.DeviceId"}, nil, nil, nil)
// 	actual := r.Get("Lorem")
// 	expected := types.DeviceId("foo2")
// 	s.Equal(expected, actual)
// 	s.Nil(r.Error())
// }
