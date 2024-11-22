package shim_provider

import (
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type provider struct {
	engine.ProviderBase
}

var _ types.ChannelProvider = (*provider)(nil)

// a shim provider to push message to the engine
// used in:
// - internal/entities/push-message/api.go to push message reveived via Rest
// - cmd/backend/main.go to push system messages like "Application started"
func NewProvider() *provider {
	return &provider{}
}
