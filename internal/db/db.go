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

	"github.com/fedulovivan/mhz19-go/internal/logger"

	"github.com/fedulovivan/mhz19-go/internal/app"
	"github.com/fedulovivan/mhz19-go/pkg/utils"
	_ "github.com/mattn/go-sqlite3"
)

var logTag = logger.MakeTag(logger.DB)

var instance *sql.DB

func DbSingleton() *sql.DB {
	if instance != nil {
		return instance
	}
	var err error
	dbabspath, err := filepath.Abs(app.Config.SqliteFilename)
	if err != nil {
		Panic(err)
	}
	if _, err := os.Stat(dbabspath); errors.Is(err, os.ErrNotExist) {
		Panic(err)
	}
	instance, err = sql.Open("sqlite3", dbabspath)
	if err != nil {
		Panic(err)
	}
	_, err = instance.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		Panic(err)
	}
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

// func Begin() (*sql.Tx, error) {
// 	return dbh.Begin()
// }

func Rollback(tx *sql.Tx) {
	if tx == nil {
		return
	}
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

func Insert(
	tx *sql.Tx,
	ctx context.Context,
	query string,
	values ...any,
) (res sql.Result, err error) {
	select {
	case <-ctx.Done():
		return
	default:
		lquery := utils.OneLineTrim(query)
		logQuery(lquery, values)
		res, err = tx.ExecContext(ctx, query, values...)
		if err != nil {
			err = fmt.Errorf(
				"got an error \"%v\" executing %v, values %v",
				err, lquery, values,
			)
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
	tx *sql.Tx,
	ctx context.Context,
	query string,
) (res int32, err error) {
	rows, err := Select(
		tx,
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

func logQuery(query string, values ...any) {
	if !app.Config.DbDebug {
		return
	}
	slog.Debug(logTag(
		fmt.Sprintf(
			"executing query %v, values %v",
			query,
			values,
		),
	))
}

func Select[T any](
	tx *sql.Tx,
	ctx context.Context,
	query string,
	scan func(rows *sql.Rows, model *T) error,
	where Where,
) (result []T, err error) {
	select {
	case <-ctx.Done():
		return
	default:
		var rows *sql.Rows
		wquery := AddWhere(query, where)
		values := PickWhereValues(where)
		lquery := utils.OneLineTrim(wquery)
		logQuery(lquery, values)
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

// https://www.reddit.com/r/golang/comments/18flz7z/comment/kcviej8/?utm_source=share&utm_medium=web3x&utm_name=web3xcss&utm_term=1&utm_content=share_button
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
