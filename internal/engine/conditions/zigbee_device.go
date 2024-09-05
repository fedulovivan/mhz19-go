package engine_conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: List
var ZigbeeDevice types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	return DeviceClass(
		mt,
		types.Args{
			"Value": types.DEVICE_CLASS_ZIGBEE_DEVICE,
		},
		e,
	) && InList(
		mt,
		types.Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
		e,
	)
}
