package engine_conditions

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// return false for nil and empty strings
// return true for the rest
var NotNil types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	c := arg_reader.NewArgReader(mt.Curr, args, nil, nil, e)
	v := c.Get("Value")
	if !c.Ok() {
		slog.Error(fmt.Sprintf("NotNil: %v", c.Error()))
		return false
	}
	switch vTyped := v.(type) {
	case string:
		return len(vTyped) > 0
	case nil:
		return false
	default:
		return true
	}
}
