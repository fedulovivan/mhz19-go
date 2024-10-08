package db

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log/slog"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/fedulovivan/mhz19-go/internal/counters"
	"github.com/fedulovivan/mhz19-go/internal/logger"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

var BaseTag = logger.NewTag(logger.DB)

type ctxkey struct{}

type ctxval struct {
	Tx  *sql.Tx
	Tag logger.Tag
}

type CtxEnhanced interface {
	context.Context
}

var instance *sql.DB

func DbSingleton() *sql.DB {
	if instance != nil {
		return instance
	}

	slog.Debug(BaseTag.F("Instance created"))

	var err error
	dbabspath, err := filepath.Abs(app.Config.SqliteFilename)
	if err != nil {
		Panic(err)
	}
	if _, err := os.Stat(dbabspath); errors.Is(err, os.ErrNotExist) {
		Panic(err)
	}
	// instance, err = sql.Open("sqlite3", fmt.Sprintf("%v?cache=shared&mode=wal", dbabspath))
	instance, err = sql.Open("sqlite3", dbabspath)
	if err != nil {
		Panic(err)
	}
	_, err = instance.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		Panic(err)
	}

	// aid for the "database is locked" issue
	// https://github.com/mattn/go-sqlite3/issues/274#issuecomment-191597862
	instance.SetMaxOpenConns(1)

	return instance
}

func Panic(err error) {
	panic("sqlite db init: " + err.Error())
}

func NewNullInt32(v int32) sql.NullInt32 {
	return sql.NullInt32{Int32: v, Valid: true}
}

func NewNullInt32FromBool(v bool) sql.NullInt32 {
	if v {
		return NewNullInt32(1)
	}
	return NewNullInt32(0)
}

func NewNullString(v string) sql.NullString {
	return sql.NullString{String: v, Valid: true}
}

func Rollback(ctx CtxEnhanced) {
	ctxpayload := ctx.Value(ctxkey{}).(ctxval)
	tx := ctxpayload.Tx
	_ = tx.Rollback()
}

// https://github.com/golang/go/issues/43507
// note(!) we should swallow "context canceled" on commit
// its expected that context is already canceled after calling g.Wait()
func Commit(tx *sql.Tx) error {
	err := tx.Commit()
	if !errors.Is(err, context.Canceled) {
		return err
	}
	return nil
}

func Exec(
	ctx CtxEnhanced,
	query string,
	values ...any,
) (res sql.Result, err error) {

	defer counters.TimeSince(time.Now(), counters.QUERIES)

	ctxpayload := ctx.Value(ctxkey{}).(ctxval)
	tx := ctxpayload.Tx
	tag := ctxpayload.Tag.WithTid("Exec")

	if app.Config.DbDebug {
		defer utils.TimeTrack(tag.F, time.Now(), "Exec")
	}

	select {
	case <-ctx.Done():
		return
	default:
		lquery := utils.OneLineTrim(query)
		// valuesj, _ := json.Marshal(values)
		logQuery(tag, lquery /* string(valuesj) */, values...)
		res, err = tx.ExecContext(ctx, query, values...)
		if err != nil {
			err = fmt.Errorf(
				"got an error \"%v\" executing %v values %v",
				err, lquery /* string(valuesj) */, values,
			)
		} else if app.Config.DbDebug {
			rows, _ := res.RowsAffected()
			slog.Debug(tag.F("Affected"), "rows", rows)
		}
		return
	}
}

type DbCount struct {
	Value int32
}

type Where = map[string]any

func AddWhere(in string, where Where) (out string) {
	var entries []string
	out = in
	for key, value := range where {
		// (!) note we cannot use "fallthrough" in type switch
		switch vtyped := value.(type) {
		case sql.NullInt32:
			if vtyped.Valid {
				entries = append(entries, fmt.Sprintf("%v = ?", key))
			}
		case sql.NullString:
			if vtyped.Valid {
				entries = append(entries, fmt.Sprintf("%v = ?", key))
			}
		}
	}
	if len(entries) > 0 {
		out = fmt.Sprintf("%v WHERE %v", in, strings.Join(entries, " AND "))
	}
	return
}

func PickWhereValues(where Where) (out []any) {
	for _, value := range where {
		switch vtyped := value.(type) {
		case sql.NullInt32:
			if vtyped.Valid {
				out = append(out, vtyped.Int32)
			}
		case sql.NullString:
			if vtyped.Valid {
				out = append(out, vtyped.String)
			}
		}
	}
	return
}

func Count(
	ctx CtxEnhanced,
	query string,
) (res int32, err error) {
	rows, err := Select(
		ctx,
		query,
		func(rows *sql.Rows, m *DbCount) error {
			return rows.Scan(&m.Value)
		},
		Where{},
	)
	if len(rows) == 1 {
		res = rows[0].Value
	} else if len(rows) > 1 {
		err = fmt.Errorf("%v query is expected to return at the most one row", query)
	}
	return
}

func logQuery(tag logger.Tag, query string, values ...any) {
	if !app.Config.DbDebug {
		return
	}
	slog.Debug(tag.F(
		"executing query %v, values %v",
		query,
		values,
	))
}

func Select[T any](
	ctx CtxEnhanced,
	query string,
	scan func(rows *sql.Rows, model *T) error,
	where Where,
) (result []T, err error) {

	defer counters.TimeSince(time.Now(), counters.QUERIES)

	ctxpayload := ctx.Value(ctxkey{}).(ctxval)
	tx := ctxpayload.Tx
	tag := ctxpayload.Tag.WithTid("Select")

	if app.Config.DbDebug {
		defer utils.TimeTrack(tag.F, time.Now(), "Select")
	}

	select {
	case <-ctx.Done():
		return
	default:
		var rows *sql.Rows
		wquery := AddWhere(query, where)
		values := PickWhereValues(where)
		lquery := utils.OneLineTrim(wquery)
		logQuery(tag, lquery, values)
		rows, err = tx.QueryContext(
			ctx,
			wquery,
			values...,
		)
		if err != nil {
			err = fmt.Errorf(
				"got an error \"%v\" for query %v, values %v",
				err, lquery, values,
			)
			return
		}
		defer rows.Close()
		for rows.Next() {
			m := new(T)
			err = scan(rows, m)
			if err != nil {
				return
			}
			result = append(result, *m)
		}
		err = rows.Err()
		return
	}
}

// initially borrowed from https://www.reddit.com/r/golang/comments/18flz7z/comment/kcviej8/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
// starts transaction
// stores tx into created context
// enhances BaseTag with unique transaction id and also stores it into created context
// commits transaction if callback returns no error
// ctx consumers are local functions Select, Exec, Rollback
func RunTx(db *sql.DB, fn func(ctx CtxEnhanced) error) error {

	defer counters.TimeSince(time.Now(), counters.TRANSACTIONS)

	tag := BaseTag.WithTid("Tx")

	if app.Config.DbDebug {
		defer utils.TimeTrack(tag.F, time.Now(), "Transaction")
	}

	tx, err := db.Begin()
	if app.Config.DbDebug {
		slog.Debug(tag.F("Transaction started"))
	}
	if err != nil {
		return err
	}
	var ctx CtxEnhanced = context.WithValue(
		context.Background(),
		ctxkey{}, ctxval{
			Tx:  tx,
			Tag: tag,
		},
	)
	defer Rollback(ctx)
	err = fn(ctx)
	if err != nil {
		return err
	}
	err = tx.Commit()
	return err
}

// if app.Config.DbDebug {
// defer func(start time.Time) {
// 	elapsed := utils.TimeTrack(tag.F, start, "Transaction")
// 	counters.Time(elapsed, counters.QUERIES)
// }(time.Now())
// }
// tx := ctx.Value(key_tx{}).(*sql.Tx)
// tag := ctx.Value(key_tag{}).(logger.Tag).AddTid("Select")
// counters.Inc(counters.QUERIES)
// if app.Config.DbDebug {
// defer func(start time.Time) {
// 	counters.Time(time.Since(start), counters.QUERIES)
// }(time.Now())
// }
// var ctx CtxEnhanced
// ctx = context.WithValue(
// 	context.Background(),
// 	key_tx{}, tx,
// )
// ctx = context.WithValue(
// 	ctx,
// 	key_tag{}, tag,
// )
// counters.Inc(counters.QUERIES)
// func WithTx(db *sql.DB, fn func(tx *sql.Tx) error) error {
//     txn, err := db.Begin()
//     if err != nil {
//         return err
//     }
//     err = fn(txn)
//     if err != nil {
//         err2 := txn.Rollback()
//         return errors.Join(err, err2)
//     }
//     return txn.Commit()
// }
// WithTx(db, func(tx *sql.Tx) error {
//     var id int
//     err := txn.QueryRow("SELECT id FROM record WHERE status = 'PENDING'").Scan(&id)
//     if err != nil {
// 	return err
//     }
//     _, err = txn.Exec("UPDATE record SET status = 'PROCESSING' WHERE id = $1", id)
//     if err != nil {
// 	return err
//     }
//     err := processRecord(id)
//     if err != nil {
// 	return err
//     }
//     _, err = txn.Exec("UPDATE record SET status = 'COMPLETED' WHERE id = $1", id)
//     if err != nil {
// 	return err
//     }
//     return nil
// })
// if app.Config.DbDebug {
// 	tag := ctx.Value(Ctxkey_tag{}).(logger.Tag)
// 	slog.Error(tag.F("Rollback"))
// }
