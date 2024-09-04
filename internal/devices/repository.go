package devices

import (
	"context"
	"database/sql"

	"github.com/fedulovivan/mhz19-go/internal/db"
)

type DbDevice struct {
	Id            int32
	NativeId      string
	DeviceClassId int32
	Name          sql.NullString
	Comments      sql.NullString
	Origin        sql.NullString
	Json          sql.NullString
}

type DevicesRepository interface {
	UpsertAll(devices []DbDevice) error
	Get(deviceId sql.NullString) ([]DbDevice, error)
}

type devicesRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) DevicesRepository {
	return devicesRepository{
		database: database,
	}
}

func (repo devicesRepository) UpsertAll(devices []DbDevice) (err error) {
	ctx := context.Background()
	tx, err := repo.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	for _, device := range devices {
		_, err = DeviceUpsert(device, ctx, tx)
		if err != nil {
			return
		}
	}
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (repo devicesRepository) Get(deviceId sql.NullString) (devices []DbDevice, err error) {
	ctx := context.Background()
	tx, err := repo.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	devices, err = DevicesSelect(ctx, tx, deviceId)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func Count(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM devices`,
	)
}

func DevicesSelect(ctx context.Context, tx *sql.Tx, nativeId sql.NullString) ([]DbDevice, error) {
	return db.Select(
		tx,
		ctx,
		`SELECT
			id,
			native_id,
			device_class_id,
			name,
			comments,
			origin,
			json
		FROM 
			devices`,
		func(rows *sql.Rows, m *DbDevice) error {
			return rows.Scan(&m.Id, &m.NativeId, &m.DeviceClassId, &m.Name, &m.Comments, &m.Origin, &m.Json)
		},
		db.Where{
			"native_id": nativeId,
		},
	)
}

func DeviceUpsert(
	device DbDevice,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO devices(native_id, device_class_id, name, comments, origin, json)
		VALUES(?,?,?,?,?,?) 
		ON CONFLICT(native_id) 
		DO UPDATE SET json = excluded.json`,
		device.NativeId, device.DeviceClassId, device.Name, device.Comments, device.Origin, device.Json,
	)
}
