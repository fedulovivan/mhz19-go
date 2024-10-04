package types

type MessagesService interface {
	Get() ([]Message, error)
	GetByDeviceId(deviceId string) ([]Message, error)
	Create(message Message) error
	CreateAll(messages []Message) error
}

type DevicesService interface {
	Get() ([]Device, error)
	GetByDeviceClass(dc DeviceClass) ([]Device, error)
	GetOne(id DeviceId) (Device, error)
	UpsertAll(devices []Device) (int64, error)
	Update(device Device) error
}

type DictItem struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
}

type DictsService interface {
	Get(DictType) ([]DictItem, error)
	All() (map[DictType][]DictItem, error)
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
	GetByDeviceId(deviceId DeviceId) (Message, error)
	OnSet() chan LdmKey
}

type StatsService interface {
	Get() (TableStats, error)
}

type RulesService interface {
	GetOne(ruleId int) (Rule, error)
	Delete(ruleId int) error
	Get() ([]Rule, error)
	Create(rule Rule) (int64, error)
	OnCreated() chan Rule
	OnDeleted() chan int
}
