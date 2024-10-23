package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: Value
var False types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (res bool, err error) {
	reader := arguments.NewReader(mt.Curr, args, nil, nil, nil, tag)
	value, err := arguments.GetTyped[bool](&reader, "Value")
	if err != nil {
		return
	}
	return Equal(
		mt,
		types.Args{
			"Left":  value,
			"Right": false,
		},
		tag,
	)
}
