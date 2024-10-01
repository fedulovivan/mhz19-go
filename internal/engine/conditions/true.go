package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var True types.CondImpl = func(mt types.MessageCompound, args types.Args, tag logger.Tag) (res bool, err error) {
	reader := arguments.NewReader(mt.Curr, args, nil, nil, nil, tag)
	value, err := arguments.GetTyped[bool](&reader, "Value")
	if err != nil {
		return
	}
	return Equal(
		mt,
		types.Args{
			"Left":  value,
			"Right": true,
		},
		tag,
	)
}
