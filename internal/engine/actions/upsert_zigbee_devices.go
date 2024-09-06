package actions

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// system action to create devices upon receiving message from zigbee2mqtt bridge
// see https://www.zigbee2mqtt.io/guide/usage/mqtt_topics_and_messages.html#zigbee2mqtt-bridge-devices
// and json example at assets/bridge-devices-message.json
var UpsertZigbeeDevices types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	devicesjson := gabs.Wrap(mm[0].Payload)
	out := make([]types.Device, 0)
	for _, d := range devicesjson.Children() {
		dtype := d.Path("type").Data().(string)
		if dtype == "Coordinator" {
			continue
		}
		out = append(out, types.Device{
			DeviceClassId: types.DEVICE_CLASS_ZIGBEE_DEVICE,
			DeviceId:      types.DeviceId(d.Path("ieee_address").Data().(string)),
			Comments:      d.Path("definition.description").Data().(string),
			Origin:        "bridge-upsert",
			Json:          d.Data(),
		})
	}
	err = e.DevicesService().UpsertAll(out)
	return
}
