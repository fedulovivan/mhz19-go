package types

import "github.com/fedulovivan/mhz19-go/internal/app"

type MessagesService interface {
	Get() ([]Message, error)
	GetByDeviceId(deviceId string) ([]Message, error)
	Create(message Message) error
}

type DevicesService interface {
	Get() ([]Device, error)
	GetByDeviceClass(dc DeviceClass) ([]Device, error)
	GetOne(id DeviceId) (Device, error)
	UpsertAll(devices []Device) error
}

type LdmKey struct {
	DeviceClass DeviceClass
	DeviceId    DeviceId
}

type LdmService interface {
	NewKey(deviceClass DeviceClass, deviceId DeviceId) LdmKey
	Get(key LdmKey) Message
	Has(key LdmKey) bool
	Set(key LdmKey, m Message)
	GetAll() []Message
	GetByDeviceId(deviceId DeviceId) Message
}

type StatsGetResult struct {
	Rules    int32 `json:"rules"`
	Devices  int32 `json:"devices"`
	Messages int32 `json:"messages"`

	EngineMessagesReceived int32 `json:"engineMessagesReceived"`
	EngineRulesMatched     int32 `json:"engineRulesMatched"`
	ApiRequests            int32 `json:"apiRequests"`
}

// TODO need better api
func (r *StatsGetResult) WithAppStats(stats *app.AppStatCounters) {
	r.EngineMessagesReceived = int32(stats.EngineMessagesReceived.Value())
	r.EngineRulesMatched = int32(stats.EngineRulesMatched.Value())
	r.ApiRequests = int32(stats.ApiRequests.Value())
}

type StatsService interface {
	Get() (StatsGetResult, error)
}

type RulesService interface {
	GetOne(ruleId int32) (Rule, error)
	Delete(ruleId int32) error
	Get() ([]Rule, error)
	Create(rule Rule) (int64, error)
	OnCreated() chan Rule
}
