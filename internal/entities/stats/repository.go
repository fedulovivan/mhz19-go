package stats

import (
	"context"
	"database/sql"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/entities/messages"
	"github.com/fedulovivan/mhz19-go/internal/entities/rules"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"golang.org/x/sync/errgroup"
)

type StatsRepository interface {
	Get() (res types.StatsGetResult, err error)
}

var _ StatsRepository = (*statsRepository)(nil)

type statsRepository struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) StatsRepository {
	return statsRepository{
		database: database,
	}
}

func (r statsRepository) Get() (
	res types.StatsGetResult,
	err error,
) {
	g, ctx := errgroup.WithContext(context.Background())
	tx, err := r.database.Begin()
	defer db.Rollback(tx)
	if err != nil {
		return
	}
	g.Go(func() (e error) { res.Rules, e = rules.CountTx(ctx, tx); return })
	g.Go(func() (e error) { res.Devices, e = devices.CountTx(ctx, tx); return })
	g.Go(func() (e error) { res.Messages, e = messages.CountTx(ctx, tx); return })
	err = g.Wait()
	if err == nil {
		err = db.Commit(tx)
	}
	return
}
