package actions

import (
	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var origin = "dnssd-upsert"

// system action to create device upon receiving dns-sd message with _ewelink._tcp service
// Args: <none>
var UpsertSonoffDevice = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	m := mm[0]
	gjson := gabs.Wrap(m.Payload)
	name := gjson.Path("Host").Data().(string)
	out := []types.Device{
		{
			DeviceId:      m.DeviceId,
			DeviceClassId: types.DEVICE_CLASS_SONOFF_DIY_PLUG,
			Name:          &name,
			Origin:        &origin,
			Json:          gjson.Data(),
		},
	}
	err = e.DevicesService().UpsertAll(out)
	return
}
