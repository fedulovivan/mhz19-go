package actions

import (
	"fmt"

	"github.com/Jeffail/gabs/v2"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: <none>
// system action to create device upon receiving dns-sd message with _ewelink._tcp service
var UpsertSonoffDevice = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.ServiceAndProviderSupplier,
	tag utils.Tag,
) (err error) {
	m := compound.Curr
	gjson := gabs.Wrap(m.Payload)
	name, ok := gjson.Path("Host").Data().(string)
	if !ok {
		err = fmt.Errorf("cannot read Host field")
		return
	}
	id, ok := gjson.Path("Id").Data().(string)
	if !ok {
		err = fmt.Errorf("cannot read Id field")
		return
	}
	origin := "dnssd-upsert"
	out := []types.Device{
		{
			DeviceId:    types.DeviceId(id),
			DeviceClass: types.DEVICE_CLASS_SONOFF_DIY_PLUG,
			Name:        &name,
			Origin:      &origin,
			Json:        gjson.Data(),
		},
	}
	_, err = e.GetDevicesService().UpsertAll(out)
	return
}
