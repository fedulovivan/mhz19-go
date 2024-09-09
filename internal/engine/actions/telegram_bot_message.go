package actions

import (
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var TelegramBotMessage types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	tpayload := types.TemplatePayload{
		Messages: mm,
	}
	areader := arguments.NewReader(&mm[0], args, mapping, &tpayload, e)
	text := areader.Get("Text")
	err = areader.Error()
	if err != nil {
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
