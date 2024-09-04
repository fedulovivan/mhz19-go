package engine_conditions

import (
	"fmt"
	"log/slog"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var InList types.CondImpl = func(mt types.MessageTuple, args types.Args, e types.Engine) bool {
	c := arg_reader.NewArgReader(mt.Curr, args, nil, nil, e)
	v := c.Get("Value")
	list := c.Get("List")
	if !c.Ok() {
		slog.Error(fmt.Sprintf("InList: %v", c.Error()))
		return false
	}
	lslice, ok := list.([]any)
	if !ok {
		panic("[]any is expected")
	}
	res := slices.Contains(lslice, v)
	return res
}
