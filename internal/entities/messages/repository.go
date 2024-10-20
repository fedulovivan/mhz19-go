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
	// Create(message DbMessage) (messageId int64, err error)
	CreateAll(messages []DbMessage) error
}

// interface guard
var _ MessagesRepository = (*messagesRepository)(nil)

type messagesRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) MessagesRepository {
	return messagesRepository{
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

// func messageInsertTx(
// 	m DbMessage,
// 	ctx context.Context,
// ) (sql.Result, error) {
// 	return db.Exec(
// 		ctx,
// 		`INSERT INTO messages(
// 			channel_type_id,
// 			device_class_id,
// 			device_id,
// 			timestamp,
// 			json
// 		) VALUES(?,?,?,?,?)`,
// 		m.ChannelTypeId,
// 		m.DeviceClassId,
// 		m.DeviceId,
// 		m.Timestamp,
// 		m.Json,
// 	)
// }

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

func (r messagesRepository) Get(deviceId sql.NullString) (messages []DbMessage, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		messages, err = messagesSelectTx(ctx, deviceId)
		return
	})
	return
}

func (r messagesRepository) CreateAll(messages []DbMessage) error {
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

// func (r messagesRepository) Create(message DbMessage) (messageId int64, err error) {
// 	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
// 		existingdevices, err := devices.DevicesSelectTx(
// 			ctx,
// 			db.NewNullString(message.DeviceId),
// 			sql.NullInt32{},
// 		)
// 		if err != nil {
// 			return
// 		}
// 		if len(existingdevices) == 0 {
// 			slog.Warn(fmt.Sprintf(
// 				"No device with class=%v id=%v in db, creating it automatically...",
// 				message.DeviceClassId,
// 				message.DeviceId,
// 			))
// 			_, err = devices.DeviceUpsertTx(devices.DbDevice{
// 				NativeId:      message.DeviceId,
// 				DeviceClassId: message.DeviceClassId,
// 				Origin:        db.NewNullString("message-autoinsert"),
// 			}, ctx)
// 			if err != nil {
// 				return
// 			}
// 		}
// 		result, err := messageInsertTx(message, ctx)
// 		if err != nil {
// 			return
// 		}
// 		messageId, err = result.LastInsertId()
// 		if err != nil {
// 			return
// 		}
// 		return
// 	})
// 	return
// }

// ctx := context.Background()
// tx, err := r.database.Begin()
// ctx = context.WithValue(ctx, db.Ctxkey_tx{}, tx)
// ctx = context.WithValue(ctx, db.Ctxkey_tag{}, db.BaseTag.WithTid("Tx"))
// defer db.Rollback(ctx)
// if err != nil {
// 	return
// }
// ctx := context.Background()
// tx, err := r.database.Begin()
// defer db.Rollback(tx)
// if err != nil {
// 	return
// }
// messages, err = messagesSelectTx(ctx, deviceId)
// if err != nil {
// 	return
// }
// err = tx.Commit()
// return
