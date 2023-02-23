package pgstore

import (
	"context"
	"fmt"
	"reflect"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"

	"github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore"
)

const (
	queryInsertNewFile = `
		INSERT INTO files (
			id,
		   	size,
		   	path%s
		) VALUES (
			$1,
			$2,
			$3%s
		)`

	queryDeleteFileByID = `
		DELETE FROM 
		    files
		WHERE
			id = $1
		RETURNING *`

	queryGetAllFiles = `
		SELECT
			id,
			size,
			path,
			created_at
		FROM
			files`
)

func TestNewPostgresStore(t *testing.T) {
	type args struct {
		dbConn *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want dbstore.DBStore
	}{
		{
			name: "success get new postgres store",
			args: args{
				dbConn: nil,
			},
			want: &postgresStore{
				dbConn: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := NewPostgresStore(tt.args.dbConn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPostgresStore() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresStore_InsertNewFile(t *testing.T) {
	type args struct {
		ctx  context.Context
		file dbstore.FileDetail
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(sqlMock sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name: "successfully inserting new file",
			args: args{
				ctx: context.Background(),
				file: dbstore.FileDetail{
					ID:   "sample-id",
					Size: 12345,
					Path: "filepath/sample-id",
				},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				query := fmt.Sprintf(queryInsertNewFile, "", "")
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
			},
			wantErr: false,
		},
		{
			name: "successfully re-inserting file",
			args: args{
				ctx: context.Background(),
				file: dbstore.FileDetail{
					ID:        "sample-id",
					Size:      12345,
					Path:      "filepath/sample-id",
					CreatedAt: time.Date(2023, 1, 1, 1, 0, 0, 0, time.UTC),
				},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				query := fmt.Sprintf(queryInsertNewFile, ", created_at", ", $4")
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(nil)
			},
			wantErr: false,
		},
		{
			name: "failed to insert record to DB",
			args: args{
				ctx: context.Background(),
				file: dbstore.FileDetail{
					ID:   "sample-id",
					Size: 12345,
					Path: "filepath/sample-id",
				},
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				query := fmt.Sprintf(queryInsertNewFile, "", "")
				sqlMock.ExpectExec(query).WillReturnResult(sqlmock.NewResult(0, 0)).WillReturnError(fmt.Errorf("some-error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("error when opening a database connection: %v\n", err)
			}
			defer mockDB.Close()
			tt.mockFunc(sqlMock)

			ps := &postgresStore{
				dbConn: sqlx.NewDb(mockDB, "postgres"),
			}
			if err := ps.InsertNewFile(tt.args.ctx, tt.args.file); (err != nil) != tt.wantErr {
				t.Errorf("InsertNewFile() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func Test_postgresStore_DeleteFileByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(sqlMock sqlmock.Sqlmock)
		want     dbstore.FileDetail
		wantErr  bool
	}{
		{
			name: "successfully delete the file",
			args: args{
				ctx: context.Background(),
				id:  "sample-id",
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "size", "path", "created_at"})
				rows.AddRow("sample-id", 123, "storage/sample-id", time.Time{})
				sqlMock.ExpectQuery(queryDeleteFileByID).WillReturnRows(rows)
			},
			want: dbstore.FileDetail{
				ID:        "sample-id",
				Size:      123,
				Path:      "storage/sample-id",
				CreatedAt: time.Time{},
			},
			wantErr: false,
		},
		{
			name: "no file deleted",
			args: args{
				ctx: context.Background(),
				id:  "sample-id",
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "size", "path", "created_at"})
				sqlMock.ExpectQuery(queryDeleteFileByID).WillReturnRows(rows)
			},
			want:    dbstore.FileDetail{},
			wantErr: true,
		},
		{
			name: "failed to delete the file",
			args: args{
				ctx: context.Background(),
				id:  "sample-id",
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(queryDeleteFileByID).WillReturnError(fmt.Errorf("some-error"))
			},
			want:    dbstore.FileDetail{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("error when opening a database connection: %v\n", err)
			}
			defer mockDB.Close()
			tt.mockFunc(sqlMock)

			ps := &postgresStore{
				dbConn: sqlx.NewDb(mockDB, "postgres"),
			}
			got, err := ps.DeleteFileByID(tt.args.ctx, tt.args.id)
			if (err != nil) != tt.wantErr {
				t.Errorf("DeleteFileByID() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("DeleteFileByID() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_postgresStore_GetAllFiles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(sqlMock sqlmock.Sqlmock)
		want     []dbstore.FileDetail
		wantErr  bool
	}{
		{
			name: "successfully get all files",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				rows := sqlmock.NewRows([]string{"id", "size", "path", "created_at"})
				rows.AddRow("sample-id-1", 111, "storage/sample-id-1", time.Time{})
				rows.AddRow("sample-id-2", 222, "storage/sample-id-2", time.Time{})
				rows.AddRow("sample-id-3", 333, "storage/sample-id-3", time.Time{})
				sqlMock.ExpectQuery(queryGetAllFiles).WillReturnRows(rows)
			},
			want: []dbstore.FileDetail{
				{
					ID:        "sample-id-1",
					Size:      111,
					Path:      "storage/sample-id-1",
					CreatedAt: time.Time{},
				},
				{
					ID:        "sample-id-2",
					Size:      222,
					Path:      "storage/sample-id-2",
					CreatedAt: time.Time{},
				},
				{
					ID:        "sample-id-3",
					Size:      333,
					Path:      "storage/sample-id-3",
					CreatedAt: time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "failed to do DB query",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectQuery(queryGetAllFiles).WillReturnError(fmt.Errorf("some-error"))
			},
			want:    []dbstore.FileDetail{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, sqlMock, err := sqlmock.New(sqlmock.QueryMatcherOption(sqlmock.QueryMatcherEqual))
			if err != nil {
				t.Errorf("error when opening a database connection: %v\n", err)
			}
			defer mockDB.Close()
			tt.mockFunc(sqlMock)

			ps := &postgresStore{
				dbConn: sqlx.NewDb(mockDB, "postgres"),
			}
			got, err := ps.GetAllFiles(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllFiles() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllFiles() got = %v, want %v", got, tt.want)
			}
		})
	}
}
