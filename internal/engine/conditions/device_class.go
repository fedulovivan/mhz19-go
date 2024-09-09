package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var DeviceClass types.CondImpl = func(mt types.MessageTuple, args types.Args) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$deviceClass",
			"Right": args["Value"],
		},
	)
}
