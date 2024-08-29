package engine

import "github.com/fedulovivan/mhz19-go/internal/logger"

type Options struct {
	logTag         logger.LogTagFn
	providers      []Provider
	rules          []Rule
	messageService MessagesService
	devicesService DevicesService
}

func NewOptions() Options {
	return Options{
		logTag: func(m string) string { return "" },
	}
}

func (o *Options) SetLogTag(f func(m string) string) {
	o.logTag = f
}
func (o *Options) SetProviders(s ...Provider) {
	o.providers = s
}
func (o *Options) SetMessagesService(s MessagesService) {
	o.messageService = s
}
func (o *Options) SetDevicesService(s DevicesService) {
	o.devicesService = s
}
func (o *Options) SetRules(rules ...Rule) {
	o.rules = rules
}
