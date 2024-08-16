package engine

type ServiceBase struct {
	Out MessageChan
}

func (s *ServiceBase) Receive() MessageChan {
	return s.Out
}

func (s *ServiceBase) Send(a ...any) {
	panic("Service.Send() should be implemented in concrete service")
}

func (s *ServiceBase) Stop() {
	// noop
}

func (s *ServiceBase) Channel() ChannelType {
	return CHANNEL_UNKNOWN
}
