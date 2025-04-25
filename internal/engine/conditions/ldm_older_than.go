package conditions

import (
	"time"

	"github.com/fedulovivan/mhz19-go/internal/arguments"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

// args: Value
// check whether last device message is older than N seconds
// notes:
// designed to work with conditions on otherDeviceId (when mt.Curr is retrieved from "last device messages" repository)
// therefore there is a possibility that Curr will be nil (see more details in MessageCompound)
var LdmOlderThan types.CondImpl = func(mt types.MessageCompound, args types.Args, tag utils.Tag) (bool, error) {
	message := mt.Curr
	if message == nil {
		return true, nil
	}
	reader := arguments.NewReader(message, args, nil, nil, nil, tag)

	// TODO use same approach as in "type Throttle struct"
	value, err := arguments.GetTyped[float64](&reader, "Value")

	if err != nil {
		return true, err
	}
	_ = value
	return time.Since(message.Timestamp) > time.Second*time.Duration(value), nil
}
