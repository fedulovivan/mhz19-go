package actions

import (
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: <none>
var RecordMessage types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.EngineAsSupplier,
	tag logger.Tag,
) error {
	slog.Debug(tag.F("Messages to save"), "len", len(compound.Queued))
	return e.MessagesService().CreateAll(compound.Queued)
}
