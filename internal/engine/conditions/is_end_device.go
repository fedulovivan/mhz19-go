package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// Args: <none>
var FromEndDevice types.CondImpl = func(mt types.MessageCompound, args types.Args, tag logger.Tag) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$fromEndDevice",
			"Right": true,
		},
		tag,
	)
}
