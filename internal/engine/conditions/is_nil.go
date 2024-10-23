package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: Value
// return false for nil and empty strings
// return true for the rest
var Nil types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (bool, error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil, tag)
	v := c.Get("Value")
	err := c.Error()
	if err != nil {
		return true, err
	}
	switch vTyped := v.(type) {
	case string:
		return len(vTyped) == 0, nil
	case nil:
		return true, nil
	default:
		return false, nil
	}
}
