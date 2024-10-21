package arguments

import (
	"testing"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/mocks"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

type ReaderSuite struct {
	suite.Suite
	tag logger.Tag
}

func (s *ReaderSuite) SetupSuite() {
}

func (s *ReaderSuite) TeardownSuite() {
}

func (s *ReaderSuite) Test10() {
	r := NewReader(nil, types.Args{"foo": 1}, nil, nil, nil, s.tag)
	actual := r.Get("foo")
	expected := 1
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test30() {
	r := NewReader(&types.Message{DeviceId: "foo3"}, types.Args{"Lorem1": "$deviceId"}, nil, nil, nil, s.tag)
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
	r := NewReader(&m, args, nil, nil, nil, s.tag)
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
	r := NewReader(&m, args, nil, nil, nil, s.tag)
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
	r := NewReader(&m, args, nil, nil, nil, s.tag)
	actual := r.Get("Lorem3")
	s.Nil(actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test70() {
	m := types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_BRIDGE}
	args := types.Args{"Lorem1": "$deviceClass"}
	r := NewReader(&m, args, nil, nil, nil, s.tag)
	actual := r.Get("Lorem1")
	s.Equal(types.DEVICE_CLASS_ZIGBEE_BRIDGE, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test80() {
	m := types.Message{}
	args := types.Args{}
	mapping := types.Mapping{}
	r := NewReader(&m, args, mapping, nil, nil, s.tag)
	s.NotNil(r)
}

func (s *ReaderSuite) Test90() {
	args := types.Args{}
	r := NewReader(nil, args, nil, nil, nil, s.tag)
	s.NotNil(r)
	r.Get("Foo")
	s.EqualError(r.Error(), "no such argument Foo")
	r.Get("$message.bar")
	s.EqualError(r.Error(), "no such argument Foo\nno such argument $message.bar")
}

func (s *ReaderSuite) Test100() {
	m := types.Message{}
	args := types.Args{
		"Foo1": "{{ range .Queued }}{{end}}",
	}
	r := NewReader(&m, args, nil, &types.TemplatePayload{Queued: []types.Message{}}, nil, s.tag)
	s.Equal("", r.Get("Foo1"))
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test101() {
	m := types.Message{}
	args := types.Args{
		"Foo2": "{{ .Test }}",
	}
	r := NewReader(&m, args, nil, &types.TemplatePayload{Queued: []types.Message{}}, nil, s.tag)
	s.Equal(r.Get("Foo2"), "{{ .Test }}")
	s.EqualError(r.Error(), `template: Foo2:1:3: executing "Foo2" at <.Test>: can't evaluate field Test in type *types.TemplatePayload`)
}

func (s *ReaderSuite) Test102() {
	m := types.Message{}
	args := types.Args{
		"Foo2": "{{ }} }}",
	}
	r := NewReader(&m, args, nil, &types.TemplatePayload{Queued: []types.Message{}}, nil, s.tag)
	s.Equal(r.Get("Foo2"), "{{ }} }}")
	s.EqualError(r.Error(), `template: Foo2:1: missing value for command`)
}

func (s *ReaderSuite) Test110() {
	m := types.Message{}
	args := types.Args{"Foo3": "lorem"}
	r := NewReader(&m, args, types.Mapping{"Foo3": {"lorem": "dolor"}}, nil, nil, s.tag)
	s.Equal("dolor", r.Get("Foo3"))
}

func (s *ReaderSuite) Test111() {
	m := types.Message{}
	args := types.Args{"Foo4": 1}
	r := NewReader(&m, args, types.Mapping{"Foo4": {"1": "2"}}, nil, nil, s.tag)
	s.Equal("2", r.Get("Foo4"))
}

func (s *ReaderSuite) Test120() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ range .Queued }}{{ .DeviceId }} - {{ deviceName .DeviceId }}{{end}}"}
	tpayload := types.TemplatePayload{Queued: []types.Message{{
		DeviceId: "10011cec96",
	}}}
	r := NewReader(&m, args, nil, &tpayload, engine, s.tag)
	s.Equal("DeviceId(10011cec96) - My perfect name", r.Get("Foo5"))
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test121() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ range .Queued }}{{ .DeviceId }} - {{ deviceName .DeviceId }}{{end}}"}
	tpayload := types.TemplatePayload{Queued: []types.Message{{
		DeviceId: "lorem111",
	}}}
	r := NewReader(&m, args, nil, &tpayload, engine, s.tag)
	s.Equal("DeviceId(lorem111) - lorem111", r.Get("Foo5"))
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test122() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ range .Queued }}device name is '{{ deviceName .DeviceId }}'{{end}}"}
	tpayload := types.TemplatePayload{Queued: []types.Message{{
		DeviceId: "nullish-device-id",
	}}}
	r := NewReader(&m, args, nil, &tpayload, engine, s.tag)
	expected := "device name is 'Device of class , with id nullish-device-id'"
	actual := r.Get("Foo5")
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test123() {
	engine := mocks.NewEngineMock()
	m := types.Message{}
	args := types.Args{"Foo5": "{{ deviceName 111 }}"}
	tpayload := types.TemplatePayload{}
	r := NewReader(&m, args, nil, &tpayload, engine, s.tag)
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
	r := NewReader(&m, args, nil, &tpayload, engine, s.tag)
	expected1 := "online"
	actual1 := r.Get("Foo5")
	s.Equal(expected1, actual1)
	expected2 := "offline"
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
	r := NewReader(&m, args, nil, &tpayload, engine, s.tag)
	expected := "333"
	actual := r.Get("Foo5")
	s.Equal(expected, actual)
	s.Nil(r.Error())
}

func (s *ReaderSuite) Test126() {
	engine := mocks.NewEngineMock()
	m := types.Message{}

	template := "{{ if gt (len .Queued) 1 }}{{ deviceName (index .Queued 0).DeviceId }}:\n{{ range .Queued }}{{ time .Timestamp }} {{ pingerStatusName .Payload.status }}\n{{ end }}{{ else }}{{ deviceName (index .Queued 0).DeviceId }} is {{ pingerStatusName (index .Queued 0).Payload.status }}{{ end }}"

	args := types.Args{"Foo6": template}

	tpayload1 := types.TemplatePayload{
		Queued: []types.Message{
			{
				DeviceId: "lorem111",
				Payload: map[string]any{
					"status": 1,
				},
			},
			{
				DeviceId: "ipsum222",
				Payload: map[string]any{
					"status": 0,
				},
			},
		},
	}
	r1 := NewReader(&m, args, nil, &tpayload1, engine, s.tag)
	actual1 := r1.Get("Foo6")
	s.Nil(r1.Error())
	s.Equal("lorem111:\n00:00:00 online\n00:00:00 offline\n", actual1)

	tpayload2 := types.TemplatePayload{
		Queued: []types.Message{
			{
				DeviceId: "lorem111",
				Payload: map[string]any{
					"status": 1,
				},
			},
		},
	}
	r2 := NewReader(&m, args, nil, &tpayload2, engine, s.tag)
	actual2 := r2.Get("Foo6")
	s.Nil(r2.Error())
	s.Equal("lorem111 is online", actual2)
}

func (s *ReaderSuite) Test130() {
	m := types.Message{DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE}
	args := types.Args{"Left": "$message.occupancy"}
	r := NewReader(&m, args, nil, nil, nil, s.tag)
	left, err := GetTyped[bool](&r, "Left")
	s.ErrorContains(err, "cannot cast <nil> to bool")
	s.False(left)
}

func TestReader(t *testing.T) {
	suite.Run(t, &ReaderSuite{
		tag: logger.NewTag(logger.ARGS),
	})
}
