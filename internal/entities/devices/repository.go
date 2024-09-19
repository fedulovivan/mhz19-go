package devices

import (
	"context"
	"database/sql"
	"strings"

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
	BuriedTimeout sql.NullInt32
}

type DevicesRepository interface {
	UpsertAll(devices []DbDevice) error
	Get(deviceId sql.NullString, deviceClass sql.NullInt32) ([]DbDevice, error)
}

var _ DevicesRepository = (*devicesRepository)(nil)

type devicesRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) DevicesRepository {
	return devicesRepository{
		database: database,
	}
}

func (r devicesRepository) UpsertAll(devices []DbDevice) (err error) {
	return db.WithTx(r.database, func(tx *sql.Tx) error {
		ctx := context.Background()
		for _, device := range devices {
			_, err = DeviceUpsertTx(device, ctx, tx)
			if err != nil {
				return err
			}
		}
		return nil
	})
}

func (r devicesRepository) Get(deviceId sql.NullString, deviceClass sql.NullInt32) (devices []DbDevice, err error) {
	ctx := context.Background()
	tx, err := r.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	devices, err = DevicesSelectTx(ctx, tx, deviceId, deviceClass)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func CountTx(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM devices`,
	)
}

func DeviceUpsertAllTx(dd []DbDevice, ctx context.Context, tx *sql.Tx) error {
	mlen := len(dd)
	cols := 6
	p := "(?,?,?,?,?,?)"
	placehoders := make([]string, mlen)
	values := make([]any, mlen*cols)
	for i, d := range dd {
		placehoders[i] = p
		values[cols*i+0] = d.NativeId
		values[cols*i+1] = d.DeviceClassId
		values[cols*i+2] = d.Name
		values[cols*i+3] = d.Comments
		values[cols*i+4] = d.Origin
		values[cols*i+5] = d.Json
	}
	_, err := db.Exec(
		tx,
		ctx,
		`INSERT INTO devices(
			native_id, 
			device_class_id, 
			name, 
			comments, 
			origin, 
			json
		)
		VALUES `+strings.Join(placehoders, ", "),
		values...,
	)
	return err
}

func DevicesSelectTx(ctx context.Context, tx *sql.Tx, nativeId sql.NullString, deviceClass sql.NullInt32) ([]DbDevice, error) {
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
			json,
			buried_timeout
		FROM 
			devices`,
		func(rows *sql.Rows, m *DbDevice) error {
			return rows.Scan(
				&m.Id,
				&m.NativeId,
				&m.DeviceClassId,
				&m.Name,
				&m.Comments,
				&m.Origin,
				&m.Json,
				&m.BuriedTimeout,
			)
		},
		db.Where{
			"native_id":       nativeId,
			"device_class_id": deviceClass,
		},
	)
}

func DeviceUpsertTx(
	device DbDevice,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Exec(
		tx,
		ctx,
		`INSERT INTO devices(native_id, device_class_id, name, comments, origin, json)
		VALUES(?,?,?,?,?,?)
		ON CONFLICT(native_id)
		DO UPDATE SET json = excluded.json`,
		device.NativeId, device.DeviceClassId, device.Name, device.Comments, device.Origin, device.Json,
	)
}
