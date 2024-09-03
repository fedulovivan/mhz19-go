package types

type MessagesService interface {
	Get() ([]Message, error)
	GetByDeviceId(deviceId string) ([]Message, error)
	Create(message Message) error
}

type DevicesService interface {
	Get() ([]Device, error)
	GetOne(id DeviceId) (Device, error)
	Upsert(devices []Device) error
}

type LdmKey string

type LdmService interface {
	MakeKey(deviceClass DeviceClass, deviceId DeviceId) LdmKey
	Get(key LdmKey) Message
	Set(key LdmKey, m Message)
	GetAll() []Message
	GetByDeviceId(deviceId DeviceId) Message
}

type StatsGetResult struct {
	Rules    int32 `json:"rules"`
	Devices  int32 `json:"devices"`
	Messages int32 `json:"messages"`
}

type StatsService interface {
	Get() (StatsGetResult, error)
}

type RulesService interface {
	GetOne(ruleId int32) (Rule, error)
	Get() ([]Rule, error)
	Create(rule Rule) (int64, error)
}
