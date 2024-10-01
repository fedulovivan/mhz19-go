package conditions

import (
	"fmt"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value, List
var InList types.CondImpl = func(mt types.MessageCompound, args types.Args, tag logger.Tag) (res bool, err error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil, tag)
	v := c.Get("Value")
	list := c.Get("List")
	err = c.Error()
	if err != nil {
		return
	}
	lslice, ok := list.([]any)
	if !ok {
		err = fmt.Errorf("[]any is expected for List")
		return
	}
	return slices.Contains(lslice, v), nil
}
