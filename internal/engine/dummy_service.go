package engine

type service struct {
	ServiceBase
}

func (s *service) Init() {
	s.Out = make(MessageChan, 100)
	s.Out <- Message{}
}

func (s *service) Send(a ...any) {
	// noop
}

func (s *service) Channel() ChannelType {
	return CHANNEL_UNKNOWN
}
