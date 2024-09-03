package stats

import (
	"context"
	"database/sql"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/devices"
	"github.com/fedulovivan/mhz19-go/internal/messages"
	"github.com/fedulovivan/mhz19-go/internal/rules"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"golang.org/x/sync/errgroup"
)

type StatsRepository interface {
	Get() (res types.StatsGetResult, err error)
}

type statsRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) StatsRepository {
	return statsRepository{
		database: database,
	}
}

func (repo statsRepository) Get() (
	res types.StatsGetResult,
	err error,
) {
	g, ctx := errgroup.WithContext(context.Background())
	tx, err := repo.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	g.Go(func() (e error) { res.Rules, e = rules.Count(ctx, tx); return })
	g.Go(func() (e error) { res.Devices, e = devices.Count(ctx, tx); return })
	g.Go(func() (e error) { res.Messages, e = messages.Count(ctx, tx); return })
	err = g.Wait()
	if err == nil {
		err = db.Commit(tx)
	}
	return
}
