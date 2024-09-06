package actions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Zigbee2MqttSetState types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	tpayload := arg_reader.TemplatePayload{
		Messages: mm,
	}
	areader := arg_reader.NewArgReader(&mm[0], args, mapping, &tpayload, e)
	deviceId := areader.Get("DeviceId")
	data := areader.Get("Data")
	if !areader.Ok() {
		err = areader.Error()
		return
	}
	topic := fmt.Sprintf("zigbee2mqtt/%v/set/state", deviceId)
	p := e.Provider(types.CHANNEL_MQTT)
	err = p.Send(topic, data)
	return
}
