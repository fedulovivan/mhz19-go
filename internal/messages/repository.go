package messages

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/devices"
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
	Get() ([]DbMessage, error)
	Create(message DbMessage) (messageId int64, err error)
}

type messagesRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) MessagesRepository {
	return messagesRepository{
		database: database,
	}
}

func messageInsert(
	m DbMessage,
	ctx context.Context,
	tx *sql.Tx,
) (sql.Result, error) {
	return db.Insert(
		tx,
		ctx,
		`INSERT INTO messages(channel_type_id,device_class_id,device_id,timestamp,json) VALUES(?,?,?,?,?)`,
		m.ChannelTypeId, m.DeviceClassId, m.DeviceId, m.Timestamp, m.Json,
	)
}

func messagesSelect(ctx context.Context, tx *sql.Tx) ([]DbMessage, error) {
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
		db.Where{},
	)
}

func Count(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM messages`,
	)
}

func (repo messagesRepository) Get() (messages []DbMessage, err error) {
	ctx := context.Background()
	tx, err := repo.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	messages, err = messagesSelect(ctx, tx)
	if err != nil {
		return
	}
	err = tx.Commit()
	return
}

func (repo messagesRepository) Create(message DbMessage) (messageId int64, err error) {
	ctx := context.Background()
	tx, err := repo.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	existingdevices, err := devices.DevicesSelect(ctx, tx, db.NewNullString(message.DeviceId))
	if err != nil {
		return
	}
	if len(existingdevices) == 0 {
		slog.Warn( /* logTag */ (fmt.Sprintf("No device with class=%v id=%v in db, creating it automatically...", message.DeviceClassId, message.DeviceId)))
		_, err = devices.DeviceUpsert(devices.DbDevice{
			NativeId:      message.DeviceId,
			DeviceClassId: message.DeviceClassId,
			Origin:        db.NewNullString("message-autoinsert"),
		}, ctx, tx)
		if err != nil {
			return
		}
	}
	result, err := messageInsert(message, ctx, tx)
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
