package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: Left, Right
var Equal types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (res bool, err error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil, tag)
	left := c.Get("Left")
	right := c.Get("Right")
	err = c.Error()
	if err != nil {
		return
	}
	return left == right, nil
}
