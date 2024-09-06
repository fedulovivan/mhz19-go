package conditions

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Equal types.CondImpl = func(mt types.MessageTuple, args types.Args) bool {
	c := arg_reader.NewArgReader(mt.Curr, args, nil, nil, nil)
	left := c.Get("Left")
	right := c.Get("Right")
	if c.Ok() {
		return left == right
	}
	slog.Error(fmt.Sprintf("Equal: %v", c.Error()))
	return false
}
