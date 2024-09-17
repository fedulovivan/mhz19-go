package engine

import "github.com/fedulovivan/mhz19-go/internal/types"

var _ types.ChannelProvider = (*ProviderBase)(nil)

type ProviderBase struct {
	Out types.MessageChan
}

func (s *ProviderBase) Messages() types.MessageChan {
	return s.Out
}

func (s *ProviderBase) Init() {
	panic("Init() must be implemented in concrete provider")
}

func (s *ProviderBase) Send(a ...any) error {
	panic("Send() must be implemented in concrete provider")
}

func (s *ProviderBase) Stop() {
	// noop
}

func (s *ProviderBase) Channel() types.ChannelType {
	panic("Channel() must be implemented in concrete provider")
}
