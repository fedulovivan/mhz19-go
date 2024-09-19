package messages

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"slices"
	"strings"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/logger"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
)

var tag = logger.NewTag(logger.DB)

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
	Create(message DbMessage) (messageId int64, err error)
	CreateAll(messages []DbMessage) error
}

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
	ctx context.Context,
	tx *sql.Tx,
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
		tx,
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

func messageInsertTx(
	m DbMessage,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Exec(
		tx,
		ctx,
		`INSERT INTO messages(
			channel_type_id,
			device_class_id,
			device_id,
			timestamp,
			json
		) VALUES(?,?,?,?,?)`,
		m.ChannelTypeId,
		m.DeviceClassId,
		m.DeviceId,
		m.Timestamp,
		m.Json,
	)
}

func messagesSelectTx(ctx context.Context, tx *sql.Tx, deviceId sql.NullString) ([]DbMessage, error) {
	return db.Select(
		tx,
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

func CountTx(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM messages`,
	)
}

func (r messagesRepository) Get(deviceId sql.NullString) (messages []DbMessage, err error) {
	ctx := context.Background()
	tx, err := r.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	messages, err = messagesSelectTx(ctx, tx, deviceId)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (r messagesRepository) CreateAll(messages []DbMessage) error {
	defer utils.TimeTrack(tag.F, time.Now(), "MessagesRepository::CreateAll()")
	ctx := context.Background()
	return db.WithTx(r.database, func(tx *sql.Tx) (err error) {
		existingDeviceIds := make([]string, 0)
		missingDeviceIds := make([]string, 0)
		missingDevices := make([]devices.DbDevice, 0)
		dd, err := devices.DevicesSelectTx(
			ctx, tx,
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
			err := devices.DeviceUpsertAllTx(missingDevices, ctx, tx)
			if err != nil {
				return err
			}
		}
		return messageInsertAllTx(messages, ctx, tx)
	})
}

func (r messagesRepository) Create(message DbMessage) (messageId int64, err error) {
	ctx := context.Background()
	tx, err := r.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	existingdevices, err := devices.DevicesSelectTx(
		ctx, tx,
		db.NewNullString(message.DeviceId),
		sql.NullInt32{},
	)
	if err != nil {
		return
	}
	if len(existingdevices) == 0 {
		slog.Warn(fmt.Sprintf(
			"No device with class=%v id=%v in db, creating it automatically...",
			message.DeviceClassId,
			message.DeviceId,
		))
		_, err = devices.DeviceUpsertTx(devices.DbDevice{
			NativeId:      message.DeviceId,
			DeviceClassId: message.DeviceClassId,
			Origin:        db.NewNullString("message-autoinsert"),
		}, ctx, tx)
		if err != nil {
			return
		}
	}
	result, err := messageInsertTx(message, ctx, tx)
	if err != nil {
		return
	}
	messageId, err = result.LastInsertId()
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}
