package migration_script

import (
	"database/sql"

	"github.com/pressly/goose"
)

const (
	migrationScriptPath = "goose/"
	postgresDialect     = "postgres"
	gooseUpCommand      = "up"
)

// MigrateUp migrate the DB to the most recent version available
func MigrateUp(db *sql.DB) error {
	if err := goose.SetDialect(postgresDialect); err != nil {
		return err
	}

	return goose.Run(gooseUpCommand, db, migrationScriptPath)
}
