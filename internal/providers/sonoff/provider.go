package sonoff_provider

import (
	"log/slog"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/hashicorp/mdns"
)

type provider struct {
	engine.ProviderBase
	ticker *time.Ticker
}

var logTag = logger.MakeTag(logger.SONOFF)

var Provider types.ChannelProvider = &provider{}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_SONOFF
}

func (p *provider) Stop() {
	if p.ticker != nil {
		slog.Debug(logTag("Stopping ticker..."))
		p.ticker.Stop()
	}
}

func (p *provider) Init() {

	p.Out = make(types.MessageChan, 100)

	entriesCh := make(chan *mdns.ServiceEntry, 4)

	go func() {
		for entry := range entriesCh {
			outMsg := types.Message{
				DeviceId:    types.DeviceId(entry.Host),
				ChannelType: p.Channel(),
				DeviceClass: types.DEVICE_CLASS_SONOFF_DIY,
				Timestamp:   time.Now(),
				Payload:     entry,
			}
			p.Out <- outMsg
		}
	}()

	var query = func() {

		// google tv
		// service := "_googlecast._tcp"

		// sonoff diy-plug devices
		// service := "_ewelink._tcp"

		// this query is sent by macos's "Discovery" tool, which makes all local devices to respond
		service := "_services._dns-sd._udp"

		slog.Debug(logTag("mdns.Query() " + service))
		err := mdns.Query(&mdns.QueryParam{
			Service:     service,
			Entries:     entriesCh,
			DisableIPv6: app.Config.MdnsDisableIPv6,
			// Timeout:     time.Second * 1,
		})
		if err != nil {
			slog.Error(logTag(err.Error()))
		}
	}

	// query for the first time
	query()

	// do periodic updates
	p.ticker = time.NewTicker(
		time.Second * 60,
	)
	go func() {
		for range p.ticker.C {
			query()
		}
	}()

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
// err := mdns.Lookup("_googlecast._tcp", entriesCh)
// err := mdns.Lookup("_ewelink._tcp", entriesCh)
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
// emit several messages to emulate db load and reproduce `got an error "database is locked" executing INSERT INTO devices`
// ticker := time.NewTicker(
// 	time.Millisecond * 5,
// )
// ticks := 10
// go func() {
// 	for i := 0; i < ticks; i++ {
// 		entriesCh <- &mdns.ServiceEntry{
// 			Host: uuid.NewString(),
// 		}
// 	}
// 	// for range ticker.C {
// 	// 	if ticks == 0 {
// 	// 		ticker.Stop()
// 	// 		return
// 	// 	}
// 	// 	entriesCh <- &mdns.ServiceEntry{
// 	// 		Host: uuid.NewString(),
// 	// 	}
// 	// 	ticks--
// 	// }
// }()
