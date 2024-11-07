package stats

import (
	"database/sql"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/entities/devices"
	"github.com/fedulovivan/mhz19-go/internal/entities/messages"
	"github.com/fedulovivan/mhz19-go/internal/entities/rules"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"golang.org/x/sync/errgroup"
)

type StatsRepository interface {
	Get() (res types.TableStats, err error)
}

var _ StatsRepository = (*repo)(nil)

type repo struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) repo {
	return repo{
		database: database,
	}
}

func (r repo) Get() (
	res types.TableStats,
	err error,
) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) error {
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (e error) { res.Rules, e = rules.CountTx(ctx); return })
		g.Go(func() (e error) { res.Actions, e = rules.CountActionsTx(ctx); return })
		g.Go(func() (e error) { res.Conds, e = rules.CountCondsTx(ctx); return })
		g.Go(func() (e error) { res.Args, e = rules.CountArgsTx(ctx); return })
		g.Go(func() (e error) { res.Mappings, e = rules.CountMappingsTx(ctx); return })
		g.Go(func() (e error) { res.Devices, e = devices.CountTx(ctx); return })
		g.Go(func() (e error) { res.Messages, e = messages.CountTx(ctx); return })
		return g.Wait()
	})
	return
}
