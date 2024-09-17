package actions

import (
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// Args: Text, BotName
var TelegramBotMessage types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	tpayload := types.TemplatePayload{
		Message:  mm[0],
		Messages: mm,
	}
	areader := arguments.NewReader(&mm[0], args, mapping, &tpayload, e)
	text := areader.Get("Text")
	botName := areader.Get("BotName")
	err = areader.Error()
	if err != nil {
		return
	}
	p := e.Provider(types.CHANNEL_TELEGRAM)
	if text == nil {
		text, _ = json.Marshal(mm[0])
	}
	return p.Send(botName, text)
}
