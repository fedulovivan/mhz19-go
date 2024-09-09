package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// return false for nil and empty strings
// return true for the rest
var NotNil types.CondImpl = func(mt types.MessageTuple, args types.Args) (res bool, err error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil)
	v := c.Get("Value")
	err = c.Error()
	if err != nil {
		return
	}
	switch vTyped := v.(type) {
	case string:
		return len(vTyped) > 0, nil
	case nil:
		return false, nil
	default:
		return true, nil
	}
}
