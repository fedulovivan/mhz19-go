package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: List
var ZigbeeDevice types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (res bool, err error) {
	classMatches, err := DeviceClass(
		mt,
		types.Args{
			"Value": types.DEVICE_CLASS_ZIGBEE_DEVICE,
		},
		tag,
	)
	if err != nil {
		return
	}
	if !classMatches {
		return false, nil
	}
	idMatches, err := InList(
		mt,
		types.Args{
			"Value": "$deviceId",
			"List":  args["List"],
		},
		tag,
	)
	if err != nil {
		return
	}
	return idMatches, nil
}
