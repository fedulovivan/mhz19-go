package rest_provider

import (
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type provider struct {
	engine.ProviderBase
}

func NewProvider() types.ChannelProvider {
	return &provider{}
}

func (p *provider) Init() {
	p.InitBase()
}
