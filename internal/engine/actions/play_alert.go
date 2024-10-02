package actions

import (
	"os/exec"

	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: <none>
var PlayAlert types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.EngineAsSupplier,
	tag logger.Tag,
) (err error) {
	_, err = exec.Command(
		"mpg123",
		"./assets/siren.mp3",
	).Output()
	// slog.Debug(tag.F(string(out)))
	return
}
