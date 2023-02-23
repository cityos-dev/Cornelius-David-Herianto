package service

import (
	"context"
	"fmt"

	"github.com/jmoiron/sqlx"
)

// Service provides mechanism to check Service's health
//
//go:generate mockgen -destination mocks/mock_service.go github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service Service
type Service interface {
	GetServiceHealth(ctx context.Context) error
}

type service struct {
	dbConn *sqlx.DB
}

// New returns a new Service instance
func New(dbConn *sqlx.DB) Service {
	return &service{
		dbConn: dbConn,
	}
}

// GetServiceHealth checks all connection-related functionality and returned the error if any
func (s *service) GetServiceHealth(ctx context.Context) error {
	// Test sql connection
	err := s.dbConn.PingContext(ctx)
	if err != nil {
		return fmt.Errorf("failed to connect to sql DB: %v", err)
	}
	return nil
}
