package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
// return false for nil and empty strings
// return true for the rest
var Nil types.CondImpl = func(mt types.MessageTuple, args types.Args) (bool, error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil)
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
