package dnssd_provider

import (
	"context"
	"fmt"
	"log/slog"
	"time"

	"github.com/brutella/dnssd"
	dnssdLogger "github.com/brutella/dnssd/log"
	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

type provider struct {
	engine.ProviderBase
}

var _ types.ChannelProvider = (*provider)(nil)

var tag = utils.NewTag(logger.DNSSD)

func NewProvider() *provider {
	return &provider{}
}

func (p *provider) Type() types.ProviderType {
	return types.PROVIDER_DNSSD
}

func (p *provider) Channel() types.ChannelType {
	return types.CHANNEL_DNS_SD
}

func (p *provider) Init() {

	p.ProviderBase.Init()

	if app.Config.DnssdDebug {
		dnssdLogger.Debug.Enable()
		dnssdLogger.Info.Enable()
	}

	ctx := context.Background()

	service := "_ewelink._tcp.local."

	addFn := func(entry dnssd.BrowseEntry) {
		if entry.Text["type"] != "diy_plug" {
			slog.Warn(tag.F("Unexpected entry data %+v", entry))
			return
		}
		payload := map[string]any{
			"Id":   entry.Text["id"],
			"Port": entry.Port,
			"Ip":   entry.IPs[0].String(),
			"Text": entry.Text,
			"Host": entry.Host,
		}
		outMsg := types.Message{
			Id:            types.MessageIdSeq.Add(1),
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
		utils.Dump("Removed", e)
	}

	go func() {
		if err := dnssd.LookupType(ctx, service, addFn, rmvFn); err != nil {
			fmt.Println(err)
			slog.Error(tag.F(err.Error()))
			counters.Inc(counters.ERRORS_ALL)
			counters.Errors.WithLabelValues(logger.MOD_DNSSD).Inc()
			return
		}
	}()

}
