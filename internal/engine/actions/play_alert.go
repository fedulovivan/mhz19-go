package actions

import (
	"os/exec"

	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: <none>
var PlayAlert types.ActionImpl = func(
	compound types.MessageCompound,
	args types.Args,
	mapping types.Mapping,
	e types.EngineAsSupplier,
	tag utils.Tag,
) (err error) {
	_, err = exec.Command(
		"mpg123",
		"./assets/siren.mp3",
	).Output()
	return
}
