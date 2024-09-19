package devices

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
)

var _ types.DevicesService = (*devicesService)(nil)

type devicesService struct {
	repository DevicesRepository
}

func (s devicesService) UpsertAll(devices []types.Device) (err error) {
	return s.repository.UpsertAll(ToDb(devices))
}

func ToDb(in []types.Device) (out []DbDevice) {
	for _, d := range in {
		mjson, err := json.Marshal(d.Json)
		dbDevice := DbDevice{
			NativeId:      string(d.DeviceId),
			DeviceClassId: int32(d.DeviceClassId),
			Json:          sql.NullString{String: string(mjson), Valid: err == nil},
		}
		if d.Name != nil {
			dbDevice.Name = db.NewNullString(*d.Name)
		}
		if d.Comments != nil {
			dbDevice.Comments = db.NewNullString(*d.Comments)
		}
		if d.Origin != nil {
			dbDevice.Origin = db.NewNullString(*d.Origin)
		}
		if d.BuriedTimeout != nil {
			dbDevice.BuriedTimeout = db.NewNullInt32(int32(d.BuriedTimeout.Duration.Seconds()))
		}
		out = append(out, dbDevice)
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
			Id:            int(d.Id),
			DeviceId:      types.DeviceId(d.NativeId),
			DeviceClassId: types.DeviceClass(d.DeviceClassId),
			Json:          payload,
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

func (s devicesService) GetByDeviceClass(dc types.DeviceClass) (devices []types.Device, err error) {
	dbdev, err := s.repository.Get(
		sql.NullString{},
		db.NewNullInt32(int32(dc)),
	)
	devices = BuildDevices(dbdev)
	return
}

func (s devicesService) Get() (devices []types.Device, err error) {
	dbdev, err := s.repository.Get(
		sql.NullString{},
		sql.NullInt32{},
	)
	if err != nil {
		return
	}
	return BuildDevices(dbdev), nil
}

func (s devicesService) GetOne(id types.DeviceId) (res types.Device, err error) {
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

func NewService(r DevicesRepository) types.DevicesService {
	return devicesService{
		repository: r,
	}
}
