package actions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: State, DeviceId
var ValveSetState types.ActionImpl = func(compound types.MessageCompound, args types.Args, mapping types.Mapping, e types.EngineAsSupplier, tag logger.Tag) (err error) {
	// tpayload := types.TemplatePayload{
	// 	Messages: mm,
	// }
	areader := arguments.NewReader(compound.Curr, args, mapping /* &tpayload */, nil, e, tag)
	deviceId := areader.Get("DeviceId")
	state := areader.Get("State")
	err = areader.Error()
	if err != nil {
		return
	}
	topic := fmt.Sprintf("/VALVE/%v/STATE/SET", deviceId)
	p := e.FindProvider(types.CHANNEL_MQTT)
	err = p.Send(topic, state)
	return
}
