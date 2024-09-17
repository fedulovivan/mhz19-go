package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// Args: <none>
var FromEndDevice types.CondImpl = func(mt types.MessageTuple, args types.Args) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$fromEndDevice",
			"Right": true,
		},
	)
}
