package postgresql

import (
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresSQLConnection() (*sqlx.DB, error) {
	db, err := sqlx.Connect("postgres", "user=postgres password=password dbname=videostore host=db sslmode=disable port=4321")
	if err != nil {
		return nil, err
	}
	return db, nil
}
