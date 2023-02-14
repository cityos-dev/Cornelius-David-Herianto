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
	healthSvc := healthSvc.New(pgConn)
	healthHandler := healthHandler.New(healthSvc)

	// files service
	filesPostgresStore := filesPGStore.NewPostgresStore(pgConn)
	filesSvc := filesSvc.New(filesPostgresStore, uuidUtils)
	filesHandler := filesHandler.New(filesSvc)

	// routes definition
	g := e.Group("/v1")
	g.GET("/health", healthHandler.GetHealth)

	g.POST("/files", filesHandler.UploadFile)
	g.GET("/files/:fileID", filesHandler.GetFileByID)
	g.GET("/files", filesHandler.GetAllFiles)
	g.DELETE("/files/:fileID", filesHandler.DeleteFileByID)

	e.Logger.Fatal(e.Start(":8080"))
}
