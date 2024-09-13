package types

import (
	"fmt"
	"time"
)

type DeviceId string

type BuriedTimeout struct {
	time.Duration
}

func (d *BuriedTimeout) MarshalJSON() ([]byte, error) {
	if d == nil || d.Duration == 0 {
		return []byte("null"), nil
	}
	return []byte(fmt.Sprintf(`"%s"`, d)), nil
}

type Device struct {
	Id            int            `json:"id,omitempty"`
	DeviceId      DeviceId       `json:"deviceId,omitempty"`
	DeviceClassId DeviceClass    `json:"deviceClassId,omitempty"`
	Name          string         `json:"name,omitempty"`
	Comments      string         `json:"comments,omitempty"`
	Origin        string         `json:"origin,omitempty"`
	Json          any            `json:"json,omitempty"`
	BuriedTimeout *BuriedTimeout `json:"buriedTimeout,omitempty"`
}
