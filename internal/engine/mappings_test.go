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

func (s *dummyProvider) Init() {
	s.Out = make(types.MessageChan, 100)
	time.Sleep(time.Millisecond * 100)
	s.Out <- types.Message{}
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
func (s *dummyldmservice) OnSet() chan types.LdmKey {
	return nil
}
func (s *dummyldmservice) Set(key types.LdmKey, m types.Message) {

}
func (s *dummyldmservice) GetAll() (res []types.Message) {
	return
}
func (s *dummyldmservice) GetByDeviceId(deviceId types.DeviceId) (out types.Message) {
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
