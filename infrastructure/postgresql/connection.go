package postgresql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresSQLConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "postgres://postgres:password@db:5432/videostore?sslmode=disable")
	if err != nil {
		return nil, err
	}
	return db, nil
}
