package dicts

import (
	"database/sql"
	"fmt"

	"github.com/fedulovivan/mhz19-go/internal/db"
	"github.com/fedulovivan/mhz19-go/internal/types"
	"golang.org/x/sync/errgroup"
)

type DbDictItem struct {
	Id   int32
	Name string
}

type DictsRepository interface {
	Get(types.DictType) ([]DbDictItem, error)
	All() (
		actions []DbDictItem,
		conditions []DbDictItem,
		channels []DbDictItem,
		deviceClasses []DbDictItem,
		err error,
	)
}

type repo struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) repo {
	return repo{
		database: database,
	}
}

func actionsSelectTx(ctx db.CtxEnhanced) ([]DbDictItem, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			name
		FROM 
			action_functions`,
		func(rows *sql.Rows, m *DbDictItem) error {
			return rows.Scan(
				&m.Id,
				&m.Name,
			)
		},
		nil,
	)
}

func conditionsSelectTx(ctx db.CtxEnhanced) ([]DbDictItem, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			name
		FROM 
			condition_functions`,
		func(rows *sql.Rows, m *DbDictItem) error {
			return rows.Scan(
				&m.Id,
				&m.Name,
			)
		},
		nil,
	)
}

func channelsSelectTx(ctx db.CtxEnhanced) ([]DbDictItem, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			name
		FROM 
			channel_types`,
		func(rows *sql.Rows, m *DbDictItem) error {
			return rows.Scan(
				&m.Id,
				&m.Name,
			)
		},
		nil,
	)
}

func deviceClassesSelectTx(ctx db.CtxEnhanced) ([]DbDictItem, error) {
	return db.Select(
		ctx,
		`SELECT
			id,
			name
		FROM 
			device_classes`,
		func(rows *sql.Rows, m *DbDictItem) error {
			return rows.Scan(
				&m.Id,
				&m.Name,
			)
		},
		nil,
	)
}

func (r repo) Get(dtype types.DictType) (res []DbDictItem, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		switch dtype {
		case types.DICT_ACTIONS:
			res, err = actionsSelectTx(ctx)
		case types.DICT_CONDITIONS:
			res, err = conditionsSelectTx(ctx)
		case types.DICT_CHANNELS:
			res, err = channelsSelectTx(ctx)
		case types.DICT_DEVICE_CLASSES:
			res, err = deviceClassesSelectTx(ctx)
		default:
			err = fmt.Errorf("no such dictionary %s", dtype)
		}
		return
	})
	return
}

func (r repo) All() (
	actions []DbDictItem,
	conditions []DbDictItem,
	channels []DbDictItem,
	deviceClasses []DbDictItem,
	err error,
) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		g, ctx := errgroup.WithContext(ctx)
		g.Go(func() (e error) { actions, e = actionsSelectTx(ctx); return })
		g.Go(func() (e error) { conditions, e = conditionsSelectTx(ctx); return })
		g.Go(func() (e error) { channels, e = channelsSelectTx(ctx); return })
		g.Go(func() (e error) { deviceClasses, e = deviceClassesSelectTx(ctx); return })
		return g.Wait()
	})
	return
}
