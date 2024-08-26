package engine

// implements [engine.Provider]
type provider struct {
	ProviderBase
}

func (s *provider) Init() {
	s.Out = make(MessageChan, 100)
	s.Out <- Message{}
}

func (s *provider) Send(a ...any) {
	// noop
}

func (s *provider) Channel() ChannelType {
	return CHANNEL_UNKNOWN
}
