package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var DeviceId types.CondImpl = func(mt types.MessageTuple, args types.Args) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$deviceId",
			"Right": args["Value"],
		},
	)
}
