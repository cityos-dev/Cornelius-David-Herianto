package dbstore

import (
	"context"
	"time"
)

// FileDetail represent the detail of the file that will be stored on database
type FileDetail struct {
	ID        string
	Size      int64
	Path      string
	CreatedAt time.Time
}

// DBStore provides file-related mechanism to interact with the database
//
//go:generate mockgen -destination mocks/mock_db_store.go github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore DBStore
type DBStore interface {
	InsertNewFile(ctx context.Context, file FileDetail) error
	DeleteFileByID(ctx context.Context, id string) (FileDetail, error)
	GetAllFiles(ctx context.Context) ([]FileDetail, error)
}
