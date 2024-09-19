package actions

import (
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// Args: <none>
var RecordMessage types.ActionImpl = func(mm []types.Message, args types.Args, mapping types.Mapping, e types.EngineAsSupplier, tag logger.Tag) (err error) {
	return e.MessagesService().Create(mm[0])
}
