package actions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Text, BotName
var TelegramBotMessage types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier, tag logger.Tag) (err error) {
	tpayload := types.TemplatePayload{
		IsFirst:  len(mm) == 1,
		Message:  mm[0],
		Messages: mm,
	}
	reader := arguments.NewReader(&mm[0], args, mapping, &tpayload, e)
	text, err := arguments.GetTyped[string](&reader, "Text")
	if err != nil {
		return
	}
	botName, err := arguments.GetTyped[string](&reader, "BotName")
	if err != nil {
		return
	}
	// text := areader.Get("Text")
	// botName := areader.Get("BotName")
	// err = areader.Error()
	// if err != nil {
	// 	return
	// }
	p := e.FindProvider(types.CHANNEL_TELEGRAM)
	// if text == nil {
	// 	text, _ = json.Marshal(mm[0])
	// }
	return p.Send(botName, text)
}
