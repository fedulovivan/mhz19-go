package engine

import (
	"sync"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type ProviderBase struct {
	sync.Mutex
	messages types.MessageChan
}

func (p *ProviderBase) Messages() types.MessageChan {
	p.Lock()
	defer p.Unlock()
	return p.messages
}

func (p *ProviderBase) Init() {
	p.Lock()
	defer p.Unlock()
	p.messages = make(types.MessageChan)
}

func (p *ProviderBase) Send(a ...any) error {
	panic("Send() must be implemented in concrete provider")
}

func (p *ProviderBase) Stop() {
	p.Lock()
	defer p.Unlock()
	close(p.messages)
}

func (p *ProviderBase) Push(m types.Message) {
	p.Lock()
	defer p.Unlock()
	p.messages <- m
}

func (s *ProviderBase) Channel() types.ChannelType {
	panic("Channel() must be implemented in concrete provider")
}

func (s *ProviderBase) Type() types.ProviderType {
	panic("Type() must be implemented in concrete provider")
}
