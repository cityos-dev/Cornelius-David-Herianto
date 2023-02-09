package main

import (
	"log"

	"github.com/labstack/echo"

	"github.com/cityos-dev/Cornelius-David-Herianto/goose/migration_script"
	"github.com/cityos-dev/Cornelius-David-Herianto/infrastructure/postgresql"
	healthHandler "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/handler"
	healthSvc "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// -- services initialization --
	pgConn, err := postgresql.NewPostgresSQLConnection()
	if err != nil {
		log.Fatalf("failed to connect to DB, err: %v", err)
	}

	// setup DB
	err = migration_script.MigrateUp(pgConn.DB)
	if err != nil {
		log.Fatalf("failed to do DB migration, err: %v", err)
	}

	// health service
	healthSvc := healthSvc.New(pgConn)
	healthHandler := healthHandler.New(healthSvc)

	// routes definition
	g := e.Group("/v1")
	g.GET("/health", healthHandler.GetHealth)

	e.Logger.Fatal(e.Start(":8080"))

}
