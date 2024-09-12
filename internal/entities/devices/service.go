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
		out = append(out, DbDevice{
			NativeId:      string(d.DeviceId),
			DeviceClassId: int32(d.DeviceClassId),
			Name:          sql.NullString{String: d.Name, Valid: len(d.Name) > 0},
			Comments:      sql.NullString{String: d.Comments, Valid: len(d.Comments) > 0},
			Origin:        sql.NullString{String: d.Origin, Valid: len(d.Origin) > 0},
			Json:          sql.NullString{String: string(mjson), Valid: err == nil},
			BuriedTimeout: db.NewNullInt32(int32(d.BuriedTimeout.Duration.Seconds())),
		})
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
			Name:          d.Name.String,
			Comments:      d.Comments.String,
			Origin:        d.Origin.String,
			Json:          payload,
		}
		if d.BuriedTimeout.Valid {
			device.BuriedTimeout = types.BuriedTimeout{
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
