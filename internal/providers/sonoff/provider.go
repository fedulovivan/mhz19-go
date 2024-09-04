package sonoff_provider

import (
	"log/slog"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/hashicorp/mdns"
)

type provider struct {
	engine.ProviderBase
}

var Provider types.ChannelProvider = &provider{}

func (s *provider) Channel() types.ChannelType {
	return types.CHANNEL_SONOFF
}

func (s *provider) Init() {

	s.Out = make(types.MessageChan, 100)

	entriesCh := make(chan *mdns.ServiceEntry, 4)
	go func() {
		for entry := range entriesCh {
			s.Out <- types.Message{
				DeviceId:    types.DeviceId(entry.Host),
				ChannelType: s.Channel(),
				DeviceClass: types.DEVICE_CLASS_SONOFF_DIY,
				Timestamp:   time.Now(),
				Payload:     entry,
			}
		}
	}()

	err := mdns.Query(&mdns.QueryParam{
		Service:     "_ewelink._tcp",
		DisableIPv6: true,
		Entries:     entriesCh,
		Timeout:     time.Second * 1,
	})
	if err != nil {
		slog.Error(err.Error())
	}

}

//
//
//
//
//
//
//
//
//
//
//
//
//
//
// Make a channel for results and start listening
// entriesCh := make(chan *mdns.ServiceEntry, 4)
// go func() {
// 	for entry := range entriesCh {
// 		fmt.Printf("Got new entry: %v\n", entry)
// 	}
// }()
// mdns.Lookup("_ewelink._tcp", entriesCh)
// close(entriesCh)
// i, err := net.InterfaceByName("en0")
// if err != nil {
// 	slog.Error(err.Error())
// }
// close(entriesCh)
// mdns.Lookup("_googlecast._tcp", entriesCh)
// mdns.Lookup("_ewelink._tcp", entriesCh)
// entriesCh := make(chan *mdns.ServiceEntry, 4)
// go func() {
// 	for entry := range entriesCh {
// 		fmt.Printf("got new entry: %v\n", entry)
// 	}
// }()
// err := mdns.Query(&mdns.QueryParam{
// 	Service:     "_ewelink._tcp",
// 	DisableIPv6: true,
// 	Entries:     entriesCh,
// 	// Service:     "_ewelink._tcp",
// 	// Timeout:     time.Second * 5,
// 	// Interface:   i,
// 	// Service:     "_ewelink._tcp.local",
// })
// if err != nil {
// 	slog.Error(err.Error())
// }
