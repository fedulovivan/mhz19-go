package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var Channel types.CondImpl = func(mt types.MessageTuple, args types.Args) bool {
	return Equal(
		mt,
		types.Args{
			"Left":  "$channelType",
			"Right": args["Value"],
		},
	)
}
