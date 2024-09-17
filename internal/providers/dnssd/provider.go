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
)

type provider struct {
	engine.ProviderBase
}

var tag = logger.NewTag(logger.DNSSD)

func NewProvider() types.ChannelProvider {
	return new(provider)
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_DNS_SD
}

func (p *provider) Init() {

	p.Out = make(types.MessageChan, 100)

	ctx := context.Background()

	service := "_ewelink._tcp.local."

	addFn := func(entry dnssd.BrowseEntry) {
		if entry.Text["type"] != "diy_plug" {
			slog.Warn(tag.F(fmt.Sprintf("Unexpected entry data %+v", entry)))
			return
		}
		payload := map[string]any{
			"Port": entry.Port,
			"Ip":   entry.IPs[0].String(),
			"Text": entry.Text,
			"Host": entry.Host,
		}
		outMsg := types.Message{
			// DeviceId:    types.DeviceId(entry.Text["id"]),
			ChannelType:   p.Channel(),
			DeviceClass:   types.DEVICE_CLASS_SONOFF_ANNOUNCE,
			Timestamp:     time.Now(),
			Payload:       payload,
			FromEndDevice: false,
		}
		p.Out <- outMsg
	}

	// just swallow "onremoved" entry with no action
	// since LookupType api does not allow nil callback
	// also looks like this feature does not work properly - cb is called when device is still "online" / "not removed"
	rmvFn := func(e dnssd.BrowseEntry) {
		// utils.Dump("Removed", e)
	}

	go func() {
		if err := dnssd.LookupType(ctx, service, addFn, rmvFn); err != nil {
			fmt.Println(err)
			slog.Error(tag.F(err.Error()))
			return
		}
	}()

}
