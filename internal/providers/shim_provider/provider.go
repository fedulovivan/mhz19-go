package shim_provider

import (
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type provider struct {
	engine.ProviderBase
}

// a shim provider to push message to the engine
// used in:
// - internal/entities/push-message/api.go to push message reveived via Rest
// - cmd/backend/main.go to push system messages like "Application started"
func NewProvider() types.ChannelProvider {
	return &provider{
		ProviderBase: engine.ProviderBase{
			MessagesChan: make(types.MessageChan /* , 100 */),
		},
	}
}

func (p *provider) Init() {
	// noop
}
