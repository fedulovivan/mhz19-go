package dnssd_provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/brutella/dnssd"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type provider struct {
	engine.ProviderBase
}

var logTag = logger.MakeTag(logger.DNSSD)

var Provider types.ChannelProvider = &provider{}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_DNS_SD
}

func (p *provider) Init() {

	p.Out = make(types.MessageChan, 100)

	ctx := context.Background()

	service := "_ewelink._tcp.local."

	addFn := func(entry dnssd.BrowseEntry) {
		if entry.Text["type"] != "diy_plug" {
			return
		}
		outMsg := types.Message{
			DeviceId:    types.DeviceId(entry.Text["id"]),
			ChannelType: p.Channel(),
			DeviceClass: types.DEVICE_CLASS_SONOFF_DIY_PLUG,
			Timestamp:   time.Now(),
			Payload:     entry,
		}
		p.Out <- outMsg
	}

	rmvFn := func(e dnssd.BrowseEntry) {
		// TODO?
		utils.Dump("Removed", e)
	}

	go func() {
		if err := dnssd.LookupType(ctx, service, addFn, rmvFn); err != nil {
			fmt.Println(err)
			slog.Error(logTag(err.Error()))
			return
		}
	}()

}

// func (p *provider) Stop() {
// if p.ticker != nil {
// 	slog.Debug(logTag("Stopping ticker..."))
// 	p.ticker.Stop()
// }
// }// utils.Dump("Added", entry)
// utils.Dump("Message", outMsg)
// slog.Error("reached")
// log.Debug.Enable()
// service := "_services._dns-sd._udp.local."
// service := "_googlecast._tcp.local."
// utils.Dump("Added", entry)
// func (p *provider) Init() {
// 	p.Out = make(types.MessageChan, 100)
// 	entriesCh := make(chan *mdns.ServiceEntry, 4)
// 	go func() {
// 		for entry := range entriesCh {
// 			outMsg := types.Message{
// 				DeviceId:    types.DeviceId(entry.Host),
// 				ChannelType: p.Channel(),
// 				DeviceClass: types.DEVICE_CLASS_SONOFF_DIY,
// 				Timestamp:   time.Now(),
// 				Payload:     entry,
// 			}
// 			p.Out <- outMsg
// 		}
// 	}()
// 	var query = func() {
// 		// google tv
// 		// service := "_googlecast._tcp"
// 		// sonoff diy-plug devices
// 		// https://sonoff.tech/diy-developer/#4
// 		// on macos `dns-sd -B _ewelink._tcp`
// 		service := "_ewelink._tcp"
// 		// this query is sent by macos's "Discovery" tool, which makes all local devices to respond
// 		// see also https://github.com/stammen/dnssd-uwp/issues/8
// 		// and https://datatracker.ietf.org/doc/html/rfc6763#section-9
// 		// on ubuntu `avahi-browse -a`
// 		// service := "_services._dns-sd._udp"
// 		slog.Debug(logTag("mdns.Query() " + service))
// 		err := mdns.Query(&mdns.QueryParam{
// 			Service:     service,
// 			Entries:     entriesCh,
// 			DisableIPv6: app.Config.MdnsDisableIPv6,
// 			// Timeout:     time.Second * 1,
// 		})
// 		if err != nil {
// 			slog.Error(logTag(err.Error()))
// 		}
// 	}
// 	// query for the first time
// 	query()
// 	// do periodic updates
// 	p.ticker = time.NewTicker(
// 		time.Second * 60,
// 	)
// 	go func() {
// 		for range p.ticker.C {
// 			query()
// 		}
// 	}()
// }
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
