package conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

// args: Value
var NotChannel types.CondImpl = func(mt types.MessageTuple, args types.Args) (bool, error) {
	res, err := Channel(mt, args)
	if err != nil {
		return false, err
	}
	return !res, nil
}
