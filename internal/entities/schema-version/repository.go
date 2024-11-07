package schema_version

import (
	"database/sql"

	"github.com/fedulovivan/mhz19-go/internal/db"
)

type DbVersion struct {
	Value int32
}

type SchemaVersionRepository interface {
	GetVersion() (int32, error)
}

type repo struct {
	database *sql.DB
}

func NewRepository(database *sql.DB) repo {
	return repo{
		database: database,
	}
}

func selectVersionTx(ctx db.CtxEnhanced) ([]DbVersion, error) {
	return db.Select(
		ctx,
		`SELECT version FROM schema_version`,
		func(rows *sql.Rows, m *DbVersion) error {
			return rows.Scan(&m.Value)
		},
		nil,
	)
}

func (r repo) GetVersion() (version int32, err error) {
	err = db.RunTx(r.database, func(ctx db.CtxEnhanced) (err error) {
		rows, err := selectVersionTx(ctx)
		if err == nil {
			version = rows[0].Value
		}
		return
	})
	return
}
