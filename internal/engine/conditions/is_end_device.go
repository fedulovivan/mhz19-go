package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// Args: <none>
var FromEndDevice types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$fromEndDevice",
			"Right": true,
		},
		tag,
	)
}
