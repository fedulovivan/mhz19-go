package actions

import (
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Text?, BotName?
var TelegramBotMessage types.ActionImpl = func(compound types.MessageCompound, args types.Args, mapping types.Mapping, e types.EngineAsSupplier, tag logger.Tag) (err error) {
	tpayload := types.TemplatePayload{
		WithPrev: compound.Prev != nil,
		Queued:   compound.Queued,
	}
	if compound.Curr != nil {
		tpayload.Message = *compound.Curr
	}
	reader := arguments.NewReader(compound.Curr, args, mapping, &tpayload, e, tag)
	var botName string
	if reader.Has("BotName") {
		botName, err = arguments.GetTyped[string](&reader, "BotName")
		if err != nil {
			return
		}
	} else {
		botName = app.Config.TelegramDefaultOutBot
	}
	var text string
	if reader.Has("Text") {
		text, err = arguments.GetTyped[string](&reader, "Text")
		if err != nil {
			return
		}
	} else {
		var mjson []byte
		mjson, err = json.Marshal(compound.Curr)
		if err != nil {
			return
		}
		text = string(mjson)
	}
	p := e.FindProvider(types.CHANNEL_TELEGRAM)
	return p.Send(botName, text)
}
