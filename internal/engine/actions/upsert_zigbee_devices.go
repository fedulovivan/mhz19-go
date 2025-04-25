package actions

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: <none>
// system action to create devices upon receiving message from zigbee2mqtt bridge
// see https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html#zigbee2mqtt-bridge-devices
// and json example at assets/bridge-devices-message.json
var UpsertZigbeeDevices types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.ServiceAndProviderSupplier,
	tag utils.Tag,
) (err error) {
	devices := make([]types.Device, 0)
	var raw []interface{}
	var ok bool
	if raw, ok = compound.Curr.Payload.([]interface{}); !ok {
		return fmt.Errorf("compound.Curr.Payload is expected to be []interface{}")
	}
	out := make([]types.ZigbeeDevice, 0)
	err = utils.MapstructureDecode(raw, &out)
	if err != nil {
		return
	}
	origin := "bridge-upsert"
	for i, device := range out {
		if device.Type != "EndDevice" {
			continue
		}
		json := raw[i]
		devices = append(devices, types.Device{
			DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE,
			DeviceId:    types.DeviceId(device.IeeeAddress),
			Comments:    &device.Definition.Description,
			Origin:      &origin,
			Json:        json,
		})
	}
	id, err := e.GetDevicesService().UpsertAll(devices)
	slog.Debug(tag.F("Created"), "LastInsertId", id)
	return
}
