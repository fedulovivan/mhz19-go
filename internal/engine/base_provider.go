package engine

import "github.com/fedulovivan/mhz19-go/internal/types"

var _ types.ChannelProvider = (*ProviderBase)(nil)

type ProviderBase struct {
	messagesChan types.MessageChan
}

func (p *ProviderBase) Messages() types.MessageChan {
	return p.messagesChan
}

func (p *ProviderBase) InitBase() {
	p.messagesChan = make(types.MessageChan, 100)
}

func (p *ProviderBase) Init() {
	panic("Init() must be implemented in concrete provider")
}

func (p *ProviderBase) Send(a ...any) error {
	panic("Send() must be implemented in concrete provider")
}

func (p *ProviderBase) Stop() {
	// noop
}

func (p *ProviderBase) Push(m types.Message) {
	p.messagesChan <- m
}

func (s *ProviderBase) Channel() types.ChannelType {
	panic("Channel() must be implemented in concrete provider")
}
