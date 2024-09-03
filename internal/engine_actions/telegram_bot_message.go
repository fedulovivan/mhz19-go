package engine_actions

import (
	"encoding/json"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var TelegramBotMessage types.ActionImpl = func(mm []types.Message, a types.Action, e types.Engine) {

	tpayload := arg_reader.TemplatePayload{
		Messages: mm,
	}
	areader := arg_reader.NewArgReader(mm[0], a.Args, a.Mapping, &tpayload)
	text := areader.Get("Text")
	if !areader.Ok() {
		slog.Error(areader.Error().Error())
		return
	}

	p := e.FindProvider(types.CHANNEL_TELEGRAM)
	if text != nil {
		p.Send(text)
	} else {
		p.Send(json.Marshal(mm[0]))
	}
}
