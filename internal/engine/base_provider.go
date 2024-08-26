package engine

type ProviderBase struct {
	Out MessageChan
}

func (s *ProviderBase) Receive() MessageChan {
	return s.Out
}

func (s *ProviderBase) Send(a ...any) {
	panic("Send() should be implemented in concrete provider")
}

func (s *ProviderBase) Stop() {
	// noop
}

func (s *ProviderBase) Channel() ChannelType {
	return CHANNEL_UNKNOWN
}
