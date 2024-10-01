package actions

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: <none>
// system action to create devices upon receiving message from zigbee2mqtt bridge
// see https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html#zigbee2mqtt-bridge-devices
// and json example at assets/bridge-devices-message.json
var UpsertZigbeeDevices types.ActionImpl = func(compound types.MessageCompound, args types.Args, mapping types.Mapping, e types.EngineAsSupplier, tag logger.Tag) (err error) {
	devicesjson := gabs.Wrap(compound.Curr.Payload)
	out := make([]types.Device, 0)
	origin := "bridge-upsert"
	for _, d := range devicesjson.Children() {
		dtype := d.Path("type").Data().(string)
		if dtype == "Coordinator" {
			continue
		}
		comments := d.Path("definition.description").Data().(string)
		out = append(out, types.Device{
			DeviceClassId: types.DEVICE_CLASS_ZIGBEE_DEVICE,
			DeviceId:      types.DeviceId(d.Path("ieee_address").Data().(string)),
			Comments:      &comments,
			Origin:        &origin,
			Json:          d.Data(),
		})
	}
	err = e.DevicesService().UpsertAll(out)
	return
}
