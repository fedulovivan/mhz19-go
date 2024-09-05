package conditions

import (
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

var conditions = types.CondImpls{
	types.COND_CHANGED:       Changed,
	types.COND_EQUAL:         Equal,
	types.COND_IN_LIST:       InList,
	types.COND_NOT_EQUAL:     NotEqual,
	types.COND_NOT_NIL:       NotNil,
	types.COND_ZIGBEE_DEVICE: ZigbeeDevice,
	types.COND_DEVICE_CLASS:  DeviceClass,
}

func Get(fn types.CondFn) (action types.CondImpl) {
	action, exist := conditions[fn]
	if !exist {
		panic(fmt.Sprintf("Condition function [%v] not yet implemented", fn))
	}
	return
}