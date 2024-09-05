package actions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var RecordMessage types.ActionImpl = func(mm []types.Message, a types.Action, e types.EngineAsSupplier) (err error) {
	err = e.MessagesService().Create(mm[0])
	return
}
