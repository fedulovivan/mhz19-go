package actions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var ValveSetState types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	tpayload := types.TemplatePayload{
		Messages: mm,
	}
	areader := arguments.NewReader(&mm[0], args, mapping, &tpayload, e)
	deviceId := areader.Get("DeviceId")
	data := areader.Get("Data")
	err = areader.Error()
	if err != nil {
		return
	}
	topic := fmt.Sprintf("/VALVE/%v/STATE/SET", deviceId)
	p := e.Provider(types.CHANNEL_MQTT)
	err = p.Send(topic, data)
	return
}
