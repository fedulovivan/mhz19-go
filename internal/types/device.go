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
	// if d == nil {
	// 	return []byte("null"), nil
	// }
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
	Id            int            `json:"id,omitempty"`
	DeviceId      DeviceId       `json:"deviceId,omitempty"`
	DeviceClassId DeviceClass    `json:"deviceClassId,omitempty"`
	Name          *string        `json:"name,omitempty"`
	Comments      *string        `json:"comments,omitempty"`
	Origin        *string        `json:"origin,omitempty"`
	Json          any            `json:"json,omitempty"`
	BuriedTimeout *BuriedTimeout `json:"buriedTimeout,omitempty"`
}
