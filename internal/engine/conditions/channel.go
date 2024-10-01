package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var Channel types.CondImpl = func(mt types.MessageCompound, args types.Args, tag logger.Tag) (bool, error) {
	return Equal(
		mt,
		types.Args{
			"Left":  "$channelType",
			"Right": args["Value"],
		},
		tag,
	)
}
