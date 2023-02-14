package main

import (
	"log"
	"os"

	"github.com/labstack/echo"

	"github.com/cityos-dev/Cornelius-David-Herianto/goose/migration_script"
	"github.com/cityos-dev/Cornelius-David-Herianto/helper/uuid"
	"github.com/cityos-dev/Cornelius-David-Herianto/infrastructure/postgresql"
	filesHandler "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/handler"
	filesSvc "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/service"
	filesPGStore "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore/pgstore"
	healthHandler "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/handler"
	healthSvc "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service"
)

func main() {
	e := echo.New()
	e.HideBanner = true

	// Database connection initialization
	pgConn, err := postgresql.NewPostgresSQLConnection(os.Getenv("POSTGRES_HOST"))
	if err != nil {
		log.Fatalf("failed to connect to DB, err: %v", err)
	}

	// setup DB
	err = migration_script.MigrateUp(pgConn.DB)
	if err != nil {
		log.Fatalf("failed to do DB migration, err: %v", err)
	}

	// -- services initialization --
	// uuid utils
	uuidUtils := uuid.NewUtils()

	// health service
	healthService := healthSvc.New(pgConn)
	healthHTTPHandler := healthHandler.New(healthService)

	// files service
	filesPostgresStore := filesPGStore.NewPostgresStore(pgConn)
	filesService := filesSvc.New(filesPostgresStore, uuidUtils)
	filesHTTPHandler := filesHandler.New(filesService)

	// routes definition
	g := e.Group("/v1")
	g.GET("/health", healthHTTPHandler.GetHealth)

	g.POST("/files", filesHTTPHandler.UploadFile)
	g.GET("/files/:fileID", filesHTTPHandler.GetFileByID)
	g.GET("/files", filesHTTPHandler.GetAllFiles)
	g.DELETE("/files/:fileID", filesHTTPHandler.DeleteFileByID)

	e.Logger.Fatal(e.Start(":8080"))
}
