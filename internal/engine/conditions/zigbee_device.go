package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: List
var ZigbeeDevice types.CondImpl = func(mt types.MessageCompound, args types.Args, tag logger.Tag) (res bool, err error) {
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
	return classMatches && idMatches, nil
}
