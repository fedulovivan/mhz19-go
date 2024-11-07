package devices

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

type DevicesRepository interface {
	UpsertAll(devices []DbDevice) (int64, error)
	Get(deviceId sql.NullString, deviceClass sql.NullInt32) ([]DbDevice, error)
	Update(device DbDevice) error
	Delete(int64) error
}

// var _ DevicesRepository = (*devicesRepository)(nil)
// var _ types.DevicesService = (*devicesService)(nil)

type service struct {
	repository DevicesRepository
}

func NewService(r DevicesRepository) service {
	return service{
		repository: r,
	}
}

func (s service) UpsertAll(devices []types.Device) (int64, error) {
	return s.repository.UpsertAll(ToDbAll(devices))
}

func ToDb(d types.Device) DbDevice {
	mjson, err := json.Marshal(d.Json)
	out := DbDevice{
		NativeId:      string(d.DeviceId),
		DeviceClassId: int32(d.DeviceClass),
		Json:          sql.NullString{String: string(mjson), Valid: err == nil},
	}
	if d.Name != nil {
		out.Name = db.NewNullString(*d.Name)
	}
	if d.Comments != nil {
		out.Comments = db.NewNullString(*d.Comments)
	}
	if d.Origin != nil {
		out.Origin = db.NewNullString(*d.Origin)
	}
	if d.BuriedTimeout != nil {
		out.BuriedTimeout = db.NewNullInt32(int32(d.BuriedTimeout.Duration.Seconds()))
	}
	return out
}

func ToDbAll(in []types.Device) (out []DbDevice) {
	for _, d := range in {
		out = append(out, ToDb(d))
	}
	return
}

func BuildDevices(in []DbDevice) (out []types.Device) {
	for _, d := range in {
		var payload map[string]any
		if d.Json.Valid {
			_ = json.Unmarshal([]byte(d.Json.String), &payload)
		}
		device := types.Device{
			Id:          int(d.Id),
			DeviceId:    types.DeviceId(d.NativeId),
			DeviceClass: types.DeviceClass(d.DeviceClassId),
			Json:        payload,
		}
		if d.Name.Valid {
			device.Name = &d.Name.String
		}
		if d.Comments.Valid {
			device.Comments = &d.Comments.String
		}
		if d.Origin.Valid {
			device.Origin = &d.Origin.String
		}
		if d.BuriedTimeout.Valid {
			device.BuriedTimeout = &types.BuriedTimeout{
				Duration: time.Duration(d.BuriedTimeout.Int32) * time.Second,
			}
		}
		out = append(out, device)
	}
	return
}

// func (s devicesService) Create(device types.Device) (int32, error) {
// 	s.repository.UpsertAll()
// 	// return s.repository.Create(ToDb(device))
// 	// s.repository.UpsertAll()
// }

func (s service) Update(device types.Device) error {
	return s.repository.Update(ToDb(device))
}

func (s service) GetByDeviceClass(dc types.DeviceClass) (devices []types.Device, err error) {
	dbdev, err := s.repository.Get(
		sql.NullString{},
		db.NewNullInt32(int32(dc)),
	)
	devices = BuildDevices(dbdev)
	return
}

func (s service) Get() (devices []types.Device, err error) {
	dbdev, err := s.repository.Get(
		sql.NullString{},
		sql.NullInt32{},
	)
	if err != nil {
		return
	}
	return BuildDevices(dbdev), nil
}

func (s service) GetOne(id types.DeviceId) (res types.Device, err error) {
	dbdev, err := s.repository.Get(
		db.NewNullString(string(id)),
		sql.NullInt32{},
	)
	if len(dbdev) == 0 {
		err = fmt.Errorf("no such device")
		return
	}
	devices := BuildDevices(dbdev)
	res = devices[0]
	return
}

func (s service) Delete(id int64) (err error) {
	err = s.repository.Delete(id)
	return
}
