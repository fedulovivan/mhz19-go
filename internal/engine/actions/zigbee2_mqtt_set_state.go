package actions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: DeviceId, State
var Zigbee2MqttSetState types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.ServiceAndProviderSupplier,
	tag utils.Tag,
) (err error) {
	areader := arguments.NewReader(compound.Curr, args, mapping, nil, e, tag)
	deviceId := areader.Get("DeviceId")
	state := areader.Get("State")
	err = areader.Error()
	if err != nil {
		return
	}
	topic := fmt.Sprintf("zigbee2mqtt/%v/set/state", deviceId)
	p := e.GetProvider(types.CHANNEL_MQTT)
	err = p.Send(topic, state)
	return
}
