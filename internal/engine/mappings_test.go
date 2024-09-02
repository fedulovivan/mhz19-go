package engine

import (
	"testing"

	"github.com/stretchr/testify/suite"
)

// type dummyProvider struct {
// 	ProviderBase
// }

// func (s *dummyProvider) Init() {
// 	s.Out = make(types.MessageChan, 100)
// }

type MappingsSuite struct {
	suite.Suite
	// engine   engine
	// provider ChannelProvider
}

func (s *MappingsSuite) SetupSuite() {
}

func (s *MappingsSuite) TeardownSuite() {
}

// func (s *MappingsSuite) Test10() {
// 	provider := &dummyProvider{}
// 	opts := NewOptions()
// 	opts.SetProviders(provider)
// 	opts.SetRules([]Rule{
// 		{
// 			Id:       1,
// 			Comments: "ut rule",
// 			Actions: []Action{
// 				{Fn: ACTION_POST_SONOFF_SWITCH_MESSAGE},
// 			},
// 		},
// 	}...)
// 	engine := NewEngine(opts)
// 	engine.Start()
// 	message := Message{
// 		ChannelType: CHANNEL_MQTT,
// 		DeviceClass: DEVICE_CLASS_ZIGBEE_DEVICE,
// 		DeviceId:    DeviceId("0x00158d0004244bda"),
// 		Timestamp:   time.Now(),
// 		Payload: map[string]any{
// 			"action": "single_left",
// 		},
// 	}
// 	go provider.Write(message)
// 	time.Sleep(time.Second * 3)
// 	engine.Stop()
// }

func TestMappings(t *testing.T) {
	suite.Run(t, new(MappingsSuite))
}
