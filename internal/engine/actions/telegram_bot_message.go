package actions

import (
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var TelegramBotMessage types.ActionImpl = func(mm []types.Message, a types.Action, e types.EngineAsSupplier) (err error) {
	tpayload := arg_reader.TemplatePayload{
		Messages: mm,
	}
	areader := arg_reader.NewArgReader(&mm[0], a.Args, a.Mapping, &tpayload, e)
	text := areader.Get("Text")
	if !areader.Ok() {
		err = areader.Error()
		return
	}
	p := e.Provider(types.CHANNEL_TELEGRAM)
	if text != nil {
		err = p.Send(text)
	} else {
		err = p.Send(json.Marshal(mm[0]))
	}
	return
}
