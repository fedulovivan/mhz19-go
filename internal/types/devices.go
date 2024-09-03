package types

type DeviceId string

type Device struct {
	Id            int         `json:"id,omitempty"`
	DeviceId      DeviceId    `json:"deviceId,omitempty"`
	DeviceClassId DeviceClass `json:"deviceClassId,omitempty"`
	Name          string      `json:"name,omitempty"`
	Comments      string      `json:"comments,omitempty"`
	Origin        string      `json:"origin,omitempty"`
	Json          any         `json:"json,omitempty"`
}

// func (d DeviceId) MarshalJSON() ([]byte, error) {
// 	return []byte(fmt.Sprintf(`"DeviceId(%s)"`, d)), nil
// }
