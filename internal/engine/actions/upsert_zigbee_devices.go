package actions

import (
	"log/slog"

	"github.com/Jeffail/gabs/v2"
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
	e types.EngineAsSupplier,
	tag utils.Tag,
) (err error) {
	devicesjson := gabs.Wrap(compound.Curr.Payload)
	out := make([]types.Device, 0)
	origin := "bridge-upsert"
	for _, d := range devicesjson.Children() {
		dtype, ok := d.Path("type").Data().(string)
		if !ok || dtype != "EndDevice" {
			continue
		}
		deviceId := d.Path("ieee_address").Data().(string)
		comments := d.Path("definition.description").Data().(string)
		out = append(out, types.Device{
			DeviceClass: types.DEVICE_CLASS_ZIGBEE_DEVICE,
			DeviceId:    types.DeviceId(deviceId),
			Comments:    &comments,
			Origin:      &origin,
			Json:        d.Data(),
		})
	}
	id, err := e.DevicesService().UpsertAll(out)
	slog.Debug(tag.F("Created"), "LastInsertId", id)
	return
}
