package types

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
	OnSet() chan LdmKey
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
