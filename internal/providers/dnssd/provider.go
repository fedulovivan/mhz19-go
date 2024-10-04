package dnssd_provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/brutella/dnssd"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type provider struct {
	engine.ProviderBase
}

var tag = logger.NewTag(logger.DNSSD)

func NewProvider() types.ChannelProvider {
	return &provider{
		ProviderBase: engine.ProviderBase{
			MessagesChan: make(types.MessageChan, 100),
		},
	}
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_DNS_SD
}

func (p *provider) Init() {

	ctx := context.Background()

	service := "_ewelink._tcp.local."

	addFn := func(entry dnssd.BrowseEntry) {
		if entry.Text["type"] != "diy_plug" {
			slog.Warn(tag.F("Unexpected entry data %+v", entry))
			return
		}
		payload := map[string]any{
			"Id":   entry.Text["id"],
			"Port": fmt.Sprintf("%v", entry.Port),
			"Ip":   entry.IPs[0].String(),
			"Text": entry.Text,
			"Host": entry.Host,
		}
		outMsg := types.Message{
			Id:            types.MessageIdSeq.Inc(),
			Timestamp:     time.Now(),
			ChannelType:   p.Channel(),
			DeviceClass:   types.DEVICE_CLASS_SONOFF_ANNOUNCE,
			Payload:       payload,
			FromEndDevice: false,
		}
		p.Push(outMsg)
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
			counters.Inc(counters.ERRORS_ALL)
			return
		}
	}()

}
