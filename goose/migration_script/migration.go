package migration_script

import (
	"database/sql"
	"github.com/pressly/goose"
)

func MigrateUp(db *sql.DB) error {
	if err := goose.SetDialect("postgres"); err != nil {
		return err
	}

	return goose.Run("up", db, "goose/")
}
