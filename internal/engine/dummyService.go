package engine

type dummyservice struct {
	out MessageChan
}

func (s *dummyservice) Receive() MessageChan {
	return s.out
}

// func (s *dummyservice) Type() ChannelType {
// 	return 0
// }

func (s *dummyservice) Init() {
	s.out = make(MessageChan, 100)
	s.out <- Message{}
}

func (s *dummyservice) Stop() {

}
