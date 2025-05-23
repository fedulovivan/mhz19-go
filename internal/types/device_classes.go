package types

import (
	"encoding/json"
	"fmt"
)

type DeviceClass byte

var _ fmt.Stringer = (*DeviceClass)(nil)
var _ json.Marshaler = (*DeviceClass)(nil)

const (
	DEVICE_CLASS_ZIGBEE_DEVICE     DeviceClass = 1
	DEVICE_CLASS_PINGER            DeviceClass = 2
	DEVICE_CLASS_VALVE             DeviceClass = 3
	DEVICE_CLASS_ZIGBEE_BRIDGE     DeviceClass = 4
	DEVICE_CLASS_BOT               DeviceClass = 5
	DEVICE_CLASS_SONOFF_DIY_PLUG   DeviceClass = 6 // why class contains certain device type?
	DEVICE_CLASS_SYSTEM            DeviceClass = 7
	DEVICE_CLASS_SONOFF_ANNOUNCE   DeviceClass = 8
	DEVICE_CLASS_ESPRESENCE_DEVICE DeviceClass = 9
)

var DEVICE_CLASS_NAMES = map[DeviceClass]string{
	DEVICE_CLASS_ZIGBEE_DEVICE:     "zigbee-device",
	DEVICE_CLASS_PINGER:            "device-pinger",
	DEVICE_CLASS_VALVE:             "valves-manipulator",
	DEVICE_CLASS_ZIGBEE_BRIDGE:     "zigbee-bridge",
	DEVICE_CLASS_BOT:               "telegram-bot",
	DEVICE_CLASS_SONOFF_DIY_PLUG:   "sonoff-diy-plug",
	DEVICE_CLASS_SYSTEM:            "system",
	DEVICE_CLASS_SONOFF_ANNOUNCE:   "sonoff-announce",
	DEVICE_CLASS_ESPRESENCE_DEVICE: "espresence-device",
}

func (dc DeviceClass) String() string {
	return DEVICE_CLASS_NAMES[dc]
}

func (dc DeviceClass) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, DEVICE_CLASS_NAMES[dc])), nil
}

func (dc *DeviceClass) UnmarshalJSON(b []byte) (err error) {
	var v string
	err = json.Unmarshal(b, &v)
	if err != nil {
		return err
	}
	for deviceClass, name := range DEVICE_CLASS_NAMES {
		if fmt.Sprintf("DeviceClass(%s)", name) == v || fmt.Sprintf("DeviceClass(%d)", deviceClass) == v {
			*dc = deviceClass
			return
		}
	}
	return fmt.Errorf("failed to unmarshal %v (type=%T) to DeviceClass", v, v)
}
