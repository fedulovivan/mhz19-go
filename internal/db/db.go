package db

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"path/filepath"

	"github.com/fedulovivan/mhz19-go/internal/app"
	_ "github.com/mattn/go-sqlite3"
)

// var dbh *sql.DB

func Init() *sql.DB {
	var err error
	dbabspath, err := filepath.Abs(app.Config.SqliteFilename)
	if err != nil {
		Panic(err)
	}
	if _, err := os.Stat(dbabspath); errors.Is(err, os.ErrNotExist) {
		Panic(err)
	}
	database, err := sql.Open("sqlite3", dbabspath)
	if err != nil {
		Panic(err)
	}
	_, err = database.Exec("PRAGMA foreign_keys=ON")
	if err != nil {
		Panic(err)
	}
	return database
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
		return tx.ExecContext(ctx, query, values...)

	}
}

func Select[T any](
	tx *sql.Tx,
	ctx context.Context,
	query string,
	scan func(rows *sql.Rows, model *T) error,
) (result []T, err error) {
	select {
	case <-ctx.Done():
		return
	default:
		var rows *sql.Rows
		rows, err = tx.QueryContext(ctx, query)
		if err != nil {
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
