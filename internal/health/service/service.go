package service

import (
	"fmt"

	"github.com/jmoiron/sqlx"
)

type Service interface {
	GetServiceHealth() error
}

type service struct {
	dbConn *sqlx.DB
}

func New(dbConn *sqlx.DB) Service {
	return &service{
		dbConn: dbConn,
	}
}

func (s *service) GetServiceHealth() error {
	// Test sql connection
	err := s.dbConn.Ping()
	if err != nil {
		return fmt.Errorf("failed to connect to sql DB: %v", err)
	}
	return nil
}
