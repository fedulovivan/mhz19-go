package actions

import (
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: <none>
var RecordMessage types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.EngineAsSupplier,
	tag utils.Tag,
) error {
	slog.Debug(tag.F("Messages to save"), "len", len(compound.Queued))
	return e.MessagesService().CreateAll(compound.Queued)
}
