package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: Value
var DeviceId types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$deviceId",
			"Right": args["Value"],
		},
		tag,
	)
}
