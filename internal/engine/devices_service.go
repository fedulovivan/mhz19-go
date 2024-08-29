package engine

import (
	"database/sql"
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/devices"
)

type Device struct {
	Id            int         `json:"id,omitempty"`
	DeviceId      DeviceId    `json:"deviceId,omitempty"`
	DeviceClassId DeviceClass `json:"deviceClassId,omitempty"`
	Name          string      `json:"name,omitempty"`
	Comments      string      `json:"comments,omitempty"`
	Origin        string      `json:"origin,omitempty"`
	Json          any         `json:"json,omitempty"`
}

type DevicesService interface {
	Get() ([]Device, error)
	Upsert(devices []Device) error
}

type devicesService struct {
	repository devices.DevicesRepository
}

func (s devicesService) Upsert(devices []Device) (err error) {
	return s.repository.UpsertDevices(ToDb(devices))
}

func ToDb(in []Device) (out []devices.DbDevice) {
	for _, d := range in {
		mjson, err := json.Marshal(d.Json)
		out = append(out, devices.DbDevice{
			NativeId:      string(d.DeviceId),
			DeviceClassId: int32(d.DeviceClassId),
			Name:          sql.NullString{String: d.Name, Valid: len(d.Name) > 0},
			Comments:      sql.NullString{String: d.Comments, Valid: len(d.Comments) > 0},
			Origin:        sql.NullString{String: d.Origin, Valid: len(d.Origin) > 0},
			Json:          sql.NullString{String: string(mjson), Valid: err == nil},
		})
	}
	return
}

func BuildDevices(in []devices.DbDevice) (out []Device) {
	for _, d := range in {
		var payload any
		if d.Json.Valid {
			_ = json.Unmarshal([]byte(d.Json.String), &payload)
		}
		out = append(out, Device{
			Id:            int(d.Id),
			DeviceId:      DeviceId(d.NativeId),
			DeviceClassId: DeviceClass(d.DeviceClassId),
			Name:          d.Name.String,
			Comments:      d.Comments.String,
			Origin:        d.Origin.String,
			Json:          payload,
		})
	}
	return
}

func (s devicesService) Get() (devices []Device, err error) {
	dbdev, err := s.repository.Get()
	if err != nil {
		return
	}
	return BuildDevices(dbdev), nil
}

func NewDevicesService(r devices.DevicesRepository) DevicesService {
	return devicesService{
		repository: r,
	}
}
