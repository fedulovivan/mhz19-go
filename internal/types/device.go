package types

import (
	"encoding/json"
	"fmt"
	"strings"
)

type DeviceId string

func (d *DeviceId) String() string {
	return fmt.Sprintf(`DeviceId(%s)`, string(*d))
}

func (d DeviceId) MarshalJSON() ([]byte, error) {
	return []byte(fmt.Sprintf(`"%s"`, d.String())), nil
}

func (a *DeviceId) UnmarshalJSON(data []byte) (err error) {
	var raw string
	err = json.Unmarshal(data, &raw)
	if err != nil {
		return
	}
	if raw == "" {
		return
	}
	if !strings.HasPrefix(raw, "DeviceId(") {
		err = fmt.Errorf(`cannot parse string "%s" into DeviceId`, raw)
		return
	}
	deviceId := raw[9 : len(raw)-1]
	*a = DeviceId(deviceId)
	return
}

type Device struct {
	Id          int         `json:"id,omitempty"`
	DeviceId    DeviceId    `json:"deviceId,omitempty"`
	DeviceClass DeviceClass `json:"deviceClass,omitempty"`
	Name        *string     `json:"name,omitempty"`
	Comments    *string     `json:"comments,omitempty"`
	Origin      *string     `json:"origin,omitempty"`
	Json        any         `json:"json,omitempty"`
	// in seconds
	// when 0 no "Have not seen" messages will be delivered fot this device
	// when NULL a default value DefaultBuriedTimeout/BURIED_TIMEOUT will be used (90m)
	// when >0 customised timeout in seconds
	BuriedTimeout *BuriedTimeout `json:"buriedTimeout,omitempty"`
}
