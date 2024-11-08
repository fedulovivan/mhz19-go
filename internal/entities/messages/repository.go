package messages

import (
	"context"
	"database/sql"
	"slices"
	"strings"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
)

type DbTemperatureMessage struct {
	Temperature sql.NullFloat64
	Timestamp   int32
}

type DbMessage struct {
	Id            int32
	ChannelTypeId int32
	DeviceClassId int32
	DeviceId      string
	Timestamp     time.Time
	Json          string
}

type MessagesRepository interface {
	Get(deviceId sql.NullString) ([]DbMessage, error)
	GetWithTemperature(deviceId sql.NullString) ([]DbTemperatureMessage, error)
	CreateAll(messages []DbMessage) error
}

// interface guard
var _ MessagesRepository = (*repo)(nil)

type repo struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) repo {
	return repo{
		database: database,
	}
}

func messageInsertAllTx(
	mm []DbMessage,
	ctx db.CtxEnhanced,
) error {
	mlen := len(mm)
	cols := 5
	p := "(?,?,?,?,?)"
	placehoders := make([]string, mlen)
	values := make([]any, mlen*cols)
	for i, m := range mm {
		placehoders[i] = p
		values[cols*i+0] = m.ChannelTypeId
		values[cols*i+1] = m.DeviceClassId
		values[cols*i+2] = m.DeviceId
		values[cols*i+3] = m.Timestamp
		values[cols*i+4] = m.Json
	}
	_, err := db.Exec(
		ctx,
		`INSERT INTO messages(
			channel_type_id,
			device_class_id,
			device_id,
			timestamp,
			json
		) VALUES `+strings.Join(placehoders, ", "),
		values...,
	)
	return err
}

func (r repo) GetWithTemperature(deviceId sql.NullString) ([]DbTemperatureMessage, error) {
	var messages []DbTemperatureMessage
	err := db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		messages, err = db.Select(
			ctx,
			`SELECT 
				DISTINCT json -> '$.temperature',
				unixepoch(timestamp) 
			FROM 
				messages
			ORDER BY 
				timestamp DESC`,
			func(rows *sql.Rows, m *DbTemperatureMessage) error {
				return rows.Scan(&m.Temperature, &m.Timestamp)
			},
			db.Where{"device_id": deviceId},
		)
		return err
	})
	return messages, err
}

func messagesSelectTx(ctx context.Context, deviceId sql.NullString) ([]DbMessage, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			channel_type_id,
			device_class_id,
			device_id,
			timestamp,
			json
		FROM 
			messages`,
		func(rows *sql.Rows, m *DbMessage) error {
			return rows.Scan(&m.Id, &m.ChannelTypeId, &m.DeviceClassId, &m.DeviceId, &m.Timestamp, &m.Json)
		},
		db.Where{"device_id": deviceId},
	)
}

func CountTx(ctx context.Context) (int32, error) {
	return db.Count(
		ctx,
		`SELECT COUNT(*) FROM messages`,
	)
}

func (r repo) Get(deviceId sql.NullString) (messages []DbMessage, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		messages, err = messagesSelectTx(ctx, deviceId)
		return
	})
	return
}

func (r repo) CreateAll(messages []DbMessage) error {
	return db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		existingDeviceIds := make([]string, 0)
		missingDeviceIds := make([]string, 0)
		missingDevices := make([]devices.DbDevice, 0)
		dd, err := devices.DevicesSelectTx(
			ctx,
			sql.NullString{},
			sql.NullInt32{},
		)
		if err != nil {
			return err
		}
		for _, d := range dd {
			existingDeviceIds = append(existingDeviceIds, d.NativeId)
		}
		for _, m := range messages {
			if !slices.Contains(existingDeviceIds, m.DeviceId) && !slices.Contains(missingDeviceIds, m.DeviceId) {
				missingDeviceIds = append(missingDeviceIds, m.DeviceId)
				missingDevices = append(missingDevices, devices.DbDevice{
					NativeId:      m.DeviceId,
					DeviceClassId: m.DeviceClassId,
					Origin:        db.NewNullString("message-autoinsert"),
				})
			}
		}
		if len(missingDevices) > 0 {
			_, err := devices.UpsertAllTx(missingDevices, ctx)
			if err != nil {
				return err
			}
		}
		return messageInsertAllTx(messages, ctx)
	})
}
