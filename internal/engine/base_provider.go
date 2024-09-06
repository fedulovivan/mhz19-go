package engine

import "github.com/fedulovivan/mhz19-go/internal/types"

type ProviderBase struct {
	Out types.MessageChan
}

func (s *ProviderBase) Messages() types.MessageChan {
	return s.Out
}

func (s *ProviderBase) Send(a ...any) error {
	panic("Send() should be implemented in concrete provider")
}

func (s *ProviderBase) Write(m types.Message) {
	s.Out <- m
}

func (s *ProviderBase) Stop() {
	// noop
}

func (s *ProviderBase) Channel() types.ChannelType {
	return types.CHANNEL_UNKNOWN
}
