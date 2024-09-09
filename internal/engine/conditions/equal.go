package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Equal types.CondImpl = func(mt types.MessageTuple, args types.Args) (res bool, err error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil)
	left := c.Get("Left")
	right := c.Get("Right")
	err = c.Error()
	if err != nil {
		return
	}
	return left == right, nil
}
