package engine_actions

import (
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var RecordMessage types.ActionImpl = func(mm []types.Message, a types.Action, e types.Engine) {
	err := e.MessagesService().Create(mm[0])
	if err != nil {
		slog.Error(logTag(err.Error()))
	}
}
