package migrations

import (
	"database/sql"

	"github.com/pkg/errors"
	"github.com/pressly/goose"
)

func init() {
	goose.AddMigration(Up_0001, Down_0001)
}

const migrationTable = "migrations"

const createSecretsTableQuery = `
	CREATE TABLE IF NOT EXISTS secrets (
		hash UUID PRIMARY KEY,
		expireAfter timestamp,
		expireAfterViews int,
		message text NOT NULL
	)
`

func Run(conn *sql.DB) error {
	err := goose.SetDialect("postgres")
	if err != nil {
		return errors.Wrap(err, "failed to set sql dialect for migration")
	}

	goose.SetTableName(migrationTable)

	return goose.Run("up", conn, "./storage/migrations")
}

func Up_0001(tx *sql.Tx) error {
	if _, err := tx.Exec(createSecretsTableQuery); err != nil {
		return errors.Wrap(err, "failed to create secrets table")
	}
	return nil
}

func Down_0001(tx *sql.Tx) error {
	if _, err := tx.Exec("DROP TABLE IF EXISTS secrets"); err != nil {
		return errors.Wrap(err, "failed to clean up table")
	}
	return nil
}
