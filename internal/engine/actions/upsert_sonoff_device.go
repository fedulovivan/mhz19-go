package actions

import (
	"fmt"

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
	out := new(types.SonoffDeviceJson)
	err = utils.MapstructureDecode(m.Payload, out)
	if err != nil {
		return
	}
	if out.ID == "" || out.Host == "" {
		err = fmt.Errorf("failed to retrieve ID or Host from device json: ID=%v, Host=%v, %+v", out.ID, out.Host, m.Payload)
		return
	}
	origin := "dnssd-upsert"
	device := []types.Device{
		{
			DeviceId:    types.DeviceId(out.ID),
			DeviceClass: types.DEVICE_CLASS_SONOFF_DIY_PLUG,
			Name:        &out.Host,
			Origin:      &origin,
			Json:        m.Payload,
		},
	}
	_, err = e.GetDevicesService().UpsertAll(device)
	return
}
