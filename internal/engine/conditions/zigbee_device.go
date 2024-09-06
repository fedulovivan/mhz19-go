package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: List
var ZigbeeDevice types.CondImpl = func(mt types.MessageTuple, args types.Args) bool {
	return DeviceClass(
		mt,
		types.Args{
			"Value": types.DEVICE_CLASS_ZIGBEE_DEVICE,
		},
	) && InList(
		mt,
		types.Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
	)
}
