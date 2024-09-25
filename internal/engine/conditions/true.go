package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var True types.CondImpl = func(mt types.MessageTuple, args types.Args) (res bool, err error) {
	reader := arguments.NewReader(mt.Curr, args, nil, nil, nil)
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
	)
}
