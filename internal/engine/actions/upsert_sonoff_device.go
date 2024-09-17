package actions

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// system action to create device upon receiving dns-sd message with _ewelink._tcp service
// Args: <none>
var UpsertSonoffDevice = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	m := mm[0]
	gjson := gabs.Wrap(m.Payload)
	name, ok := gjson.Path("Host").Data().(string)
	if !ok {
		err = fmt.Errorf("cannot read Host field")
		return
	}
	origin := "dnssd-upsert"
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
