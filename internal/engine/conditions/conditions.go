package conditions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var conditions = types.CondImpls{
	types.COND_CHANGED:         Changed,
	types.COND_EQUAL:           Equal,
	types.COND_IN_LIST:         InList,
	types.COND_IS_NIL:          IsNil,
	types.COND_ZIGBEE_DEVICE:   ZigbeeDevice,
	types.COND_DEVICE_CLASS:    DeviceClass,
	types.COND_СHANNEL:         Channel,
	types.COND_FROM_END_DEVICE: FromEndDevice,
	types.COND_TRUE:            True,
	types.COND_FALSE:           False,
}

func Get(fn types.CondFn) (action types.CondImpl) {
	action, exist := conditions[fn]
	if !exist {
		panic(fmt.Sprintf("Condition function %d not yet implemented", fn))
	}
	return
}
