package pgstore

import (
	"context"
	"time"

	"github.com/jmoiron/sqlx"

	"github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore"
)

type postgresStore struct {
	dbConn *sqlx.DB
}

func NewPostgresStore(dbConn *sqlx.DB) dbstore.DBStore {
	return &postgresStore{
		dbConn: dbConn,
	}
}

type fileDetail struct {
	ID        string    `db:"id,omitempty"`
	Name      string    `db:"name,omitempty"`
	Size      int64     `db:"size,omitempty"`
	Path      string    `db:"path,omitempty"`
	CreatedAt time.Time `db:"created_at"`
}

func (ps *postgresStore) InsertNewFile(ctx context.Context, file dbstore.FileDetail) error {
	query := `
		INSERT INTO files (
			id,
		   	size,
		   	path
		) VALUES (
			:id,
			:size,
			:path
		)`

	internalFile := mapFileDetail(file)
	_, err := ps.dbConn.NamedExecContext(ctx, query, &internalFile)
	if err != nil {
		return err
	}
	return nil
}

func (ps *postgresStore) GetFileByID(ctx context.Context, id string) (dbstore.FileDetail, error) {
	//TODO implement me
	panic("implement me")
}

func (ps *postgresStore) GetAllFiles(ctx context.Context) ([]dbstore.FileDetail, error) {
	//TODO implement me
	panic("implement me")
}

func mapFileDetail(file dbstore.FileDetail) fileDetail {
	return fileDetail{
		ID:        file.ID,
		Size:      file.Size,
		Path:      file.Path,
		CreatedAt: file.CreatedAt,
	}
}
