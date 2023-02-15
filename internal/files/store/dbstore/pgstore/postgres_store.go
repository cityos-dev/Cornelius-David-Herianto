package pgstore

import (
	"context"
	"database/sql"
	"fmt"
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
	Size      int64     `db:"size,omitempty"`
	Path      string    `db:"path,omitempty"`
	CreatedAt time.Time `db:"created_at"`
}

func (ps *postgresStore) InsertNewFile(ctx context.Context, file dbstore.FileDetail) error {
	query := `
		INSERT INTO files (
			id,
		   	size,
		   	path%s
		) VALUES (
			:id,
			:size,
			:path%s
		)`

	if !file.CreatedAt.IsZero() {
		query = fmt.Sprintf(query, ", created_at", ", :created_at")
	} else {
		query = fmt.Sprintf(query, "", "")
	}

	internalFile := mapFileDetail(file)
	_, err := ps.dbConn.NamedExecContext(ctx, query, &internalFile)
	if err != nil {
		return err
	}
	return nil
}

func (ps *postgresStore) DeleteFileByID(ctx context.Context, id string) (dbstore.FileDetail, error) {
	query := `
		DELETE FROM 
		    files
		WHERE
			id = $1
		RETURNING *`

	var files []fileDetail
	err := ps.dbConn.SelectContext(ctx, &files, query, id)
	if err != nil {
		return dbstore.FileDetail{}, err
	}

	if len(files) == 0 {
		return dbstore.FileDetail{}, sql.ErrNoRows
	}

	return reverseMapFileDetail(files[0]), nil
}

func (ps *postgresStore) GetAllFiles(ctx context.Context) ([]dbstore.FileDetail, error) {
	query := `
		SELECT
			id,
			size,
			path,
			created_at
		FROM
			files`

	var files []fileDetail
	err := ps.dbConn.SelectContext(ctx, &files, query)
	if err != nil {
		return []dbstore.FileDetail{}, err
	}
	var result []dbstore.FileDetail
	for _, file := range files {
		result = append(result, reverseMapFileDetail(file))
	}
	return result, nil
}

func mapFileDetail(file dbstore.FileDetail) fileDetail {
	return fileDetail{
		ID:        file.ID,
		Size:      file.Size,
		Path:      file.Path,
		CreatedAt: file.CreatedAt,
	}
}

func reverseMapFileDetail(file fileDetail) dbstore.FileDetail {
	return dbstore.FileDetail{
		ID:        file.ID,
		Size:      file.Size,
		Path:      file.Path,
		CreatedAt: file.CreatedAt,
	}
}
