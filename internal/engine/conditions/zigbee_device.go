package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: List
var ZigbeeDevice types.CondImpl = func(mt types.MessageTuple, args types.Args) (res bool, err error) {
	classMatches, err := DeviceClass(
		mt,
		types.Args{
			"Value": types.DEVICE_CLASS_ZIGBEE_DEVICE,
		},
	)
	if err != nil {
		return
	}
	idMatches, err := InList(
		mt,
		types.Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
	)
	if err != nil {
		return
	}
	return classMatches && idMatches, nil
}
