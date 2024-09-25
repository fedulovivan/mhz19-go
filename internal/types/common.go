package types

const DEVICE_ID_FOR_THE_BURIED_DEVICES_PROVIDER_MESSAGE DeviceId = DeviceId("buried-device-id")
const DEVICE_ID_FOR_THE_REST_PROVIDER_MESSAGE DeviceId = DeviceId("rest-device-id")

type DictType string

const (
	DICT_ACTIONS        DictType = "actions"
	DICT_CONDITIONS     DictType = "conditions"
	DICT_DEVICE_CLASSES DictType = "device-classes"
	DICT_CHANNELS       DictType = "channels"
)
