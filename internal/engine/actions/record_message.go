package actions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var RecordMessage types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier) (err error) {
	return e.MessagesService().Create(mm[0])
}
