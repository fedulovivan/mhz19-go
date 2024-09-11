package types

import (
	"encoding/json"
	"fmt"
)

type DeviceClass byte

var _ fmt.Stringer = (*DeviceClass)(nil)
var _ json.Marshaler = (*DeviceClass)(nil)

// var _ json.Unmarshaler = (*DeviceClass)(nil)

const (
	DEVICE_CLASS_UNKNOWN         DeviceClass = 0
	DEVICE_CLASS_ZIGBEE_DEVICE   DeviceClass = 1
	DEVICE_CLASS_PINGER          DeviceClass = 2
	DEVICE_CLASS_VALVE           DeviceClass = 3
	DEVICE_CLASS_ZIGBEE_BRIDGE   DeviceClass = 4
	DEVICE_CLASS_BOT             DeviceClass = 5
	DEVICE_CLASS_SONOFF_DIY_PLUG DeviceClass = 6
)

var DEVICE_CLASS_NAMES = map[DeviceClass]string{
	DEVICE_CLASS_UNKNOWN:         "<unknown>",
	DEVICE_CLASS_ZIGBEE_DEVICE:   "zigbee-device",
	DEVICE_CLASS_PINGER:          "device-pinger",
	DEVICE_CLASS_VALVE:           "valve-manipulator",
	DEVICE_CLASS_ZIGBEE_BRIDGE:   "zigbee-bridge",
	DEVICE_CLASS_BOT:             "telegram-bot",
	DEVICE_CLASS_SONOFF_DIY_PLUG: "sonoff-diy-plug",
}

func (dc DeviceClass) String() string {
	return fmt.Sprintf("%v (id=%d)", DEVICE_CLASS_NAMES[dc], dc)
}

func (dc DeviceClass) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"DeviceClass(%v)"`, DEVICE_CLASS_NAMES[dc])), nil
}
