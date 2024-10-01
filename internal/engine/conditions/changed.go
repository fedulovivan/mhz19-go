package conditions

import (
	"errors"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var Changed types.CondImpl = func(mt types.MessageCompound, args types.Args, tag logger.Tag) (res bool, err error) {
	if mt.Prev == nil && mt.Curr != nil {
		return true, nil
	}
	cCurr := arguments.NewReader(mt.Curr, args, nil, nil, nil, tag)
	cPrev := arguments.NewReader(mt.Prev, args, nil, nil, nil, tag)
	vCurr := cCurr.Get("Value")
	vPrev := cPrev.Get("Value")
	err = errors.Join(cCurr.Error(), cPrev.Error())
	if err != nil {
		return
	}
	return vCurr != vPrev, nil
}
