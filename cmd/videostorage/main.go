package main

import (
	"github.com/cityos-dev/Cornelius-David-Herianto/infrastructure/postgresql"
	"github.com/labstack/echo"
	"log"

	health_handler "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/handler"
	health_svc "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// -- services initialization --
	pgConn, err := postgresql.NewPostgresSQLConnection()
	if err != nil {
		log.Fatalf("failed to connect to DB, err: %v", err)
	}

	// health service
	healthSvc := health_svc.New(pgConn)
	healthHandler := health_handler.New(healthSvc)

	// routes definition
	g := e.Group("/v1")
	g.GET("/health", healthHandler.GetHealth)

	e.Logger.Fatal(e.Start(":8080"))

}
