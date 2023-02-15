package dbstore

import (
	"context"
	"time"
)

type FileDetail struct {
	ID        string
	Size      int64
	Path      string
	CreatedAt time.Time
}

type DBStore interface {
	InsertNewFile(ctx context.Context, file FileDetail) error
	DeleteFileByID(ctx context.Context, id string) (FileDetail, error)
	GetAllFiles(ctx context.Context) ([]FileDetail, error)
}
