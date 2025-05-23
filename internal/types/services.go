package types

type MessagesService interface {
	Get() ([]Message, error)
	GetByDeviceId(DeviceId) ([]Message, error)
	GetWithTemperature(DeviceId) ([]TemperatureMessage, error)
	Create(Message) error
	CreateAll([]Message) error
}

type DevicesService interface {
	Get() ([]Device, error)
	GetByDeviceClass(dc DeviceClass) ([]Device, error)
	GetOne(id DeviceId) (Device, error)
	UpsertAll(devices []Device) (int64, error)
	UpdateName(device Device) error
	UpdateBuriedTimeout(device Device) error
	Delete(int64) error
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

// service for the "last device messages" repository
// interface replicates LdmRepository
type LdmService interface {
	NewKey(deviceClass DeviceClass, deviceId DeviceId) LdmKey
	Get(key LdmKey) Message
	Has(key LdmKey) bool
	Set(key LdmKey, m Message)
	GetAll() []Message
	GetByDeviceId(deviceId DeviceId) (Message, error)
	OnSet() <-chan LdmKey
}

type StatsService interface {
	Get() (TableStats, error)
}

type RulesService interface {
	GetOne(ruleId int) (Rule, error)
	Delete(ruleId int) error
	Get() ([]Rule, error)
	Create(rule Rule) (int64, error)
	OnCreated() <-chan Rule
	OnDeleted() <-chan int
}

type ServiceSupplier interface {
	GetDevicesService() DevicesService
	GetMessagesService() MessagesService
	GetLdmService() LdmService
}

type ServiceAndProviderSupplier interface {
	ServiceSupplier
	GetProvider(pt ProviderType) ChannelProvider
}
