package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var NotEqual types.CondImpl = func(mt types.MessageTuple, args types.Args) (bool, error) {
	res, err := Equal(mt, args)
	if err != nil {
		return false, err
	}
	return !res, nil
}
