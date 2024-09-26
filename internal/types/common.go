package types

const DEVICE_ID_FOR_THE_BURIED_DEVICES_PROVIDER_MESSAGE = DeviceId(
	"device-id-for-the-buried-devices-provider-message",
)
const DEVICE_ID_FOR_THE_REST_PROVIDER_MESSAGE = DeviceId(
	"device-id-for-the-rest-provider-message",
)
const DEVICE_ID_FOR_THE_APPLICATION_MESSAGE = DeviceId(
	"device-id-for-the-application-message",
)

type DictType string

const (
	DICT_ACTIONS        DictType = "actions"
	DICT_CONDITIONS     DictType = "conditions"
	DICT_DEVICE_CLASSES DictType = "device-classes"
	DICT_CHANNELS       DictType = "channels"
)
