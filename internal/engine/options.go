package engine

import (
	"github.com/fedulovivan/mhz19-go/internal/devices"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/messages"
)

type Options struct {
	logTag         logger.LogTagFn
	providers      []ChannelProvider
	rules          []Rule
	messageService messages.MessagesService
	devicesService devices.DevicesService
}

func NewOptions() Options {
	return Options{
		logTag: func(m string) string { return "" },
	}
}

func (o *Options) SetLogTag(f func(m string) string) {
	o.logTag = f
}
func (o *Options) SetProviders(s ...ChannelProvider) {
	o.providers = s
}
func (o *Options) SetMessagesService(s messages.MessagesService) {
	o.messageService = s
}
func (o *Options) SetDevicesService(s devices.DevicesService) {
	o.devicesService = s
}
func (o *Options) SetRules(rules ...Rule) {
	o.rules = rules
}
