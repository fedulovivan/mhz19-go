package engine

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type options struct {
	logTag         types.LogTagFn
	providers      []types.ChannelProvider
	rules          []types.Rule
	messageService types.MessagesService
	devicesService types.DevicesService
	ldmService     types.LdmService
}

func NewOptions() types.EngineOptions {
	return &options{
		logTag: func(m string) string { return "" },
	}
}

func (o *options) SetLogTag(f func(m string) string) {
	o.logTag = f
}
func (o *options) SetProviders(s ...types.ChannelProvider) {
	o.providers = s
}
func (o *options) SetMessagesService(s types.MessagesService) {
	o.messageService = s
}
func (o *options) SetDevicesService(s types.DevicesService) {
	o.devicesService = s
}
func (o *options) SetLdmService(r types.LdmService) {
	o.ldmService = r
}
func (o *options) SetRules(rules ...types.Rule) {
	o.rules = rules
}

func (o *options) LogTag() types.LogTagFn {
	return o.logTag
}
func (o *options) Providers() []types.ChannelProvider {
	return o.providers
}
func (o *options) MessagesService() types.MessagesService {
	return o.messageService
}
func (o *options) DevicesService() types.DevicesService {
	return o.devicesService
}
func (o *options) LdmService() types.LdmService {
	return o.ldmService
}
func (o *options) Rules() []types.Rule {
	return o.rules
}
