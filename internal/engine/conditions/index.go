package engine_conditions

import (
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var Conditions = types.CondImpls{
	types.COND_CHANGED:       Changed,
	types.COND_EQUAL:         Equal,
	types.COND_IN_LIST:       InList,
	types.COND_NOT_EQUAL:     NotEqual,
	types.COND_NOT_NIL:       NotNil,
	types.COND_ZIGBEE_DEVICE: ZigbeeDevice,
	types.COND_DEVICE_CLASS:  DeviceClass,
}
