package sonoff_provider

import (
	"github.com/fedulovivan/mhz19-go/internal/engine"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type provider struct {
	engine.ProviderBase
}

var Provider engine.ChannelProvider = &provider{}

func (s *provider) Channel() types.ChannelType {
	return types.CHANNEL_SONOFF
}

func (s *provider) Send(a ...any) {
}

func (s *provider) Stop() {
}

func (s *provider) SendNewMessage(text string, chatId int64) {
}

func (s *provider) Init() {
	s.Out = make(types.MessageChan, 100)
}
