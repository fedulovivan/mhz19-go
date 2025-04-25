package actions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var MQTT_TOPIC_BY_DEVICE_CLASS = map[types.DeviceClass]string{
	types.DEVICE_CLASS_ZIGBEE_DEVICE: "zigbee2mqtt/%v/set/state",
	types.DEVICE_CLASS_VALVE:         "valves-manipulator/%v/cmd",
}

// args: DeviceId, State
var MqttSetState types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.ServiceAndProviderSupplier,
	tag utils.Tag,
) (err error) {
	reader := arguments.NewReader(
		compound.Curr, args, mapping, nil, e, tag,
	)
	deviceId, err := arguments.GetTyped[types.DeviceId](&reader, "DeviceId")
	if err != nil {
		return
	}
	state, err := arguments.GetTyped[string](&reader, "State")
	if err != nil {
		return
	}
	device, err := e.GetDevicesService().GetOne(deviceId)
	if err != nil {
		return
	}
	topic, err := getStateTopicByDeviceClass(device.DeviceClass, deviceId)
	if err != nil {
		return
	}
	p := e.GetProvider(types.PROVIDER_MQTT)
	err = p.Send(topic, state)
	return
}

func getStateTopicByDeviceClass(
	dc types.DeviceClass,
	deviceId types.DeviceId,
) (string, error) {
	if pattern, ok := MQTT_TOPIC_BY_DEVICE_CLASS[dc]; ok {
		return fmt.Sprintf(pattern, deviceId), nil
	}
	return "", fmt.Errorf(
		"cannot find mqtt topic by device class %s, check MQTT_TOPIC_BY_DEVICE_CLASS",
		dc,
	)
}
