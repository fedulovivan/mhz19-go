package stats

import (
	"context"
	"database/sql"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"golang.org/x/sync/errgroup"
)

type StatsRepository interface {
	Get() (res GetResult, err error)
}

type statsRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) StatsRepository {
	return statsRepository{
		database: database,
	}
}

func rulesCount(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM rules`,
	)
}

func devicesCount(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM devices`,
	)
}

func messagesCount(ctx context.Context, tx *sql.Tx) (int32, error) {
	return db.Count(
		tx,
		ctx,
		`SELECT COUNT(*) FROM messages`,
	)
}

func (repo statsRepository) Get() (
	res GetResult,
	err error,
) {
	// defer utils.TimeTrack(logTag, time.Now(), "repo:Get")
	g, ctx := errgroup.WithContext(context.Background())
	tx, err := repo.database.Begin()
	if err != nil {
		return
	}
	g.Go(func() (e error) { res.Rules, e = rulesCount(ctx, tx); return })
	g.Go(func() (e error) { res.Devices, e = devicesCount(ctx, tx); return })
	g.Go(func() (e error) { res.Messages, e = messagesCount(ctx, tx); return })
	err = g.Wait()
	if err == nil {
		err = db.Commit(tx)
	}
	return
}
