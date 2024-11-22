package engine

import (
	"testing"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/stretchr/testify/suite"
)

var _ types.LdmService = (*dummyldmservice)(nil)

type dummyProvider struct {
	ProviderBase
}

func (p *dummyProvider) Init() {
	p.ProviderBase.Init()
	time.Sleep(time.Millisecond * 100)
	p.Push(types.Message{
		Id:        types.MessageIdSeq.Add(1),
		Timestamp: time.Now(),
	})
}

type MappingsSuite struct {
	suite.Suite
}

func (s *MappingsSuite) SetupSuite() {
}

func (s *MappingsSuite) TeardownSuite() {
}

type dummyldmservice struct {
}

func (s *dummyldmservice) NewKey(deviceClass types.DeviceClass, deviceId types.DeviceId) (res types.LdmKey) {
	return
}
func (s *dummyldmservice) Get(key types.LdmKey) (out types.Message) {
	return
}
func (s *dummyldmservice) Has(key types.LdmKey) bool {
	return false
}
func (s *dummyldmservice) OnSet() <-chan types.LdmKey {
	return nil
}
func (s *dummyldmservice) Set(key types.LdmKey, m types.Message) {

}
func (s *dummyldmservice) GetAll() (res []types.Message) {
	return
}
func (s *dummyldmservice) GetByDeviceId(deviceId types.DeviceId) (out types.Message, err error) {
	return
}

func (s *MappingsSuite) Test10() {
	engine := NewEngine()
	engine.SetProviders(&dummyProvider{})
	engine.SetLdmService(&dummyldmservice{})
	engine.Start()
	time.Sleep(time.Second * 1)
}

func TestMappings(t *testing.T) {
	suite.Run(t, new(MappingsSuite))
}
