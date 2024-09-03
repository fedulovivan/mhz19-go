package devices

import (
	"database/sql"
	"encoding/json"

	"github.com/fedulovivan/mhz19-go/internal/types"
)

type devicesService struct {
	repository DevicesRepository
}

func (s devicesService) Upsert(devices []types.Device) (err error) {
	return s.repository.UpsertDevices(ToDb(devices))
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
		})
	}
	return
}

func BuildDevices(in []DbDevice) (out []types.Device) {
	for _, d := range in {
		var payload any
		if d.Json.Valid {
			_ = json.Unmarshal([]byte(d.Json.String), &payload)
		}
		out = append(out, types.Device{
			Id:            int(d.Id),
			DeviceId:      types.DeviceId(d.NativeId),
			DeviceClassId: types.DeviceClass(d.DeviceClassId),
			Name:          d.Name.String,
			Comments:      d.Comments.String,
			Origin:        d.Origin.String,
			Json:          payload,
		})
	}
	return
}

func (s devicesService) Get() (devices []types.Device, err error) {
	dbdev, err := s.repository.Get()
	if err != nil {
		return
	}
	return BuildDevices(dbdev), nil
}

func (s devicesService) GetOne(id types.DeviceId) (res types.Device, err error) {
	return
}

func NewService(r DevicesRepository) types.DevicesService {
	return devicesService{
		repository: r,
	}
}
