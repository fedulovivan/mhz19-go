package conditions

import (
	"fmt"
	"log/slog"

	"github.com/fedulovivan/mhz19-go/internal/arg_reader"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Changed types.CondImpl = func(mt types.MessageTuple, args types.Args) bool {
	if mt.Prev == nil && mt.Curr != nil {
		return true
	}
	cCurr := arg_reader.NewArgReader(mt.Curr, args, nil, nil, nil)
	rPrev := arg_reader.NewArgReader(mt.Prev, args, nil, nil, nil)
	vCurr := cCurr.Get("Value")
	vPrev := rPrev.Get("Value")
	if cCurr.Ok() && rPrev.Ok() {
		return vCurr != vPrev
	}
	slog.Error(fmt.Sprintf("Changed: curr %v, prev %v", cCurr.Error(), rPrev.Error()))
	return false
}
