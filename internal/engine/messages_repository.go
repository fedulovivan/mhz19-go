package engine

import (
	"context"
	"database/sql"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/db"
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

func deviceInsert() {}

func (repo messagesRepository) Create(message DbMessage) (messageId int64, err error) {
	ctx := context.Background()
	tx, err := repo.database.Begin()
	if err != nil {
		return
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
