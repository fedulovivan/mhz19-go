package conditions

import (
	"fmt"
	"slices"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var InList types.CondImpl = func(mt types.MessageTuple, args types.Args) (res bool, err error) {
	c := arguments.NewReader(mt.Curr, args, nil, nil, nil)
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
