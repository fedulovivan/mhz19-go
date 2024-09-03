package engine_actions

import (
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var RecordMessage types.ActionImpl = func(mm []types.Message, a types.Action, e types.Engine) {
	err := e.GetOptions().MessagesService().Create(mm[0])
	if err != nil {
		slog.Error(e.GetOptions().LogTag()(err.Error()))
	}
}
