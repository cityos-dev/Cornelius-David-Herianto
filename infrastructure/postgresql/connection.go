package postgresql

import (
	"fmt"

	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func NewPostgresSQLConnection(host string) (*sqlx.DB, error) {
	if host == "" {
		host = "localhost"
	}
	dsn := fmt.Sprintf("user=postgres password=password dbname=videostore host=%s sslmode=disable", host)
	db, err := sqlx.Connect("postgres", dsn)
	if err != nil {
		return nil, err
	}
	return db, nil
}
