package devices

import (
	"database/sql"
	"fmt"
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
	UpsertAll(devices []DbDevice) (int64, error)
	Get(deviceId sql.NullString, deviceClass sql.NullInt32) ([]DbDevice, error)
	Update(device DbDevice) error
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

func (r devicesRepository) Update(device DbDevice) error {
	return db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		_, err = UpdateTx(ctx, device)
		return
	})
}

func (r devicesRepository) UpsertAll(devices []DbDevice) (id int64, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		res, err := UpsertAllTx(devices, ctx)
		if err != nil {
			return
		}
		id, _ = res.LastInsertId()
		return
	})
	return
}

func (r devicesRepository) Get(deviceId sql.NullString, deviceClass sql.NullInt32) (devices []DbDevice, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		devices, err = DevicesSelectTx(ctx, deviceId, deviceClass)
		return
	})
	return
}

func CountTx(ctx db.CtxEnhanced) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM devices`,
	)
}

func UpsertAllTx(devices []DbDevice, ctx db.CtxEnhanced) (sql.Result, error) {
	mlen := len(devices)
	cols := 6
	p := "(?,?,?,?,?,?)"
	placehoders := make([]string, mlen)
	values := make([]any, mlen*cols)
	for i, d := range devices {
		placehoders[i] = p
		values[cols*i+0] = d.NativeId
		values[cols*i+1] = d.DeviceClassId
		values[cols*i+2] = d.Name
		values[cols*i+3] = d.Comments
		values[cols*i+4] = d.Origin
		values[cols*i+5] = d.Json
	}
	return db.Exec(
		ctx,
		fmt.Sprintf(
			`INSERT INTO devices(
				native_id, 
				device_class_id, 
				name, 
				comments, 
				origin, 
				json
			)
			VALUES %s ON CONFLICT(native_id) DO UPDATE SET json = excluded.json`,
			strings.Join(placehoders, ", "),
		),
		values...,
	)
}

func DevicesSelectTx(ctx db.CtxEnhanced, nativeId sql.NullString, deviceClass sql.NullInt32) ([]DbDevice, error) {
	return db.Select(
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

func UpdateTx(
	ctx db.CtxEnhanced,
	device DbDevice,
) (sql.Result, error) {
	return db.Exec(
		ctx,
		`UPDATE devices SET name = ? WHERE native_id = ?`,
		device.Name,
		device.NativeId,
	)
}

// for _, device := range devices {
// 	_, err = DeviceUpsertTx(device, ctx)
// 	if err != nil {
// 		return err
// 	}
// }
// return nil

// func (r devicesRepository) Create(device DbDevice) (deviceId int32, err error) {
// 	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
// 		// insertTx
// 		// _, err = updateTx(ctx, device)
// 		// return
// 	})
// }

// func DeviceUpsertTx(
// 	device DbDevice,
// 	ctx db.CtxEnhanced,
// ) (sql.Result, error) {
// 	return db.Exec(
// 		ctx,
// 		`INSERT INTO devices(native_id, device_class_id, name, comments, origin, json)
// 		VALUES(?,?,?,?,?,?)
// 		ON CONFLICT(native_id)
// 		DO UPDATE SET json = excluded.json`,
// 		device.NativeId, device.DeviceClassId, device.Name, device.Comments, device.Origin, device.Json,
// 	)
// }
