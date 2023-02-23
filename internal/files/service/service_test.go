package service

import (
	"context"
	"fmt"
	"mime/multipart"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/lib/pq"

	"github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore"
	dbStoreMocks "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		dbStore dbstore.DBStore
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "successfully get new Service",
			args: args{
				dbStore: nil,
			},
			want: service{
				dbStore: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.dbStore); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

type mockMultipartFile struct {
	reader *strings.Reader
}

func (m mockMultipartFile) Read(p []byte) (n int, err error) {
	return m.reader.Read(p)
}

func (m mockMultipartFile) ReadAt(p []byte, off int64) (n int, err error) {
	return m.ReadAt(p, off)
}

func (m mockMultipartFile) Seek(offset int64, whence int) (int64, error) {
	return m.reader.Seek(offset, whence)
}

func (m mockMultipartFile) Close() error {
	return nil
}

func Test_service_UploadFile(t *testing.T) {
	type args struct {
		ctx      context.Context
		file     multipart.File
		host     string
		filename string
		size     int64
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockDBStore *dbStoreMocks.MockDBStore)
		want     string
		wantErr  bool
	}{
		{
			name: "successfully upload a file",
			args: args{
				ctx: context.Background(),
				file: mockMultipartFile{
					reader: strings.NewReader("sample string"),
				},
				host:     "localhost",
				filename: "test.mp4",
				size:     123,
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().InsertNewFile(context.Background(), dbstore.FileDetail{
					ID:   "test.mp4",
					Size: 123,
					Path: "localhost/v1/files/test.mp4",
				}).Return(nil)
			},
			want:    "localhost/v1/files/test.mp4",
			wantErr: false,
		},
		{
			name: "unsupported file type",
			args: args{
				ctx: context.Background(),
				file: mockMultipartFile{
					reader: strings.NewReader("sample string"),
				},
				host:     "localhost",
				filename: "text.txt",
				size:     123,
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "failed to insert to DB",
			args: args{
				ctx: context.Background(),
				file: mockMultipartFile{
					reader: strings.NewReader("sample string"),
				},
				host:     "localhost",
				filename: "test.mp4",
				size:     123,
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().InsertNewFile(context.Background(), dbstore.FileDetail{
					ID:   "test.mp4",
					Size: 123,
					Path: "localhost/v1/files/test.mp4",
				}).Return(fmt.Errorf("some-error"))
			},
			want:    "",
			wantErr: true,
		},
		{
			name: "duplicate entry",
			args: args{
				ctx: context.Background(),
				file: mockMultipartFile{
					reader: strings.NewReader("sample string"),
				},
				host:     "localhost",
				filename: "test.mp4",
				size:     123,
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().InsertNewFile(context.Background(), dbstore.FileDetail{
					ID:   "test.mp4",
					Size: 123,
					Path: "localhost/v1/files/test.mp4",
				}).Return(&pq.Error{Code: "23505"})
			},
			want:    "",
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDBStore := dbStoreMocks.NewMockDBStore(ctrl)

			tt.mockFunc(mockDBStore)

			s := service{
				dbStore: mockDBStore,
			}
			got, err := s.UploadFile(tt.args.ctx, tt.args.file, tt.args.host, tt.args.filename, tt.args.size)
			if (err != nil) != tt.wantErr {
				t.Errorf("UploadFile() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("UploadFile() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_service_GetAllFiles(t *testing.T) {
	type args struct {
		ctx context.Context
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockDBStore *dbStoreMocks.MockDBStore)
		want     []FileInfo
		wantErr  bool
	}{
		{
			name: "successfully get all files",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().GetAllFiles(context.Background()).Return([]dbstore.FileDetail{
					{
						ID:        "file-1.mp4",
						Size:      1111,
						Path:      "path/to/file-1.mp4",
						CreatedAt: time.Time{},
					},
					{
						ID:        "file-2.mp4",
						Size:      2222,
						Path:      "path/to/file-2.mp4",
						CreatedAt: time.Time{},
					},
					{
						ID:        "file-3.mp4",
						Size:      3333,
						Path:      "path/to/file-3.mp4",
						CreatedAt: time.Time{},
					},
				}, nil)
			},
			want: []FileInfo{
				{
					FileID:    "file-1.mp4",
					Size:      1111,
					Name:      "file-1.mp4",
					CreatedAt: time.Time{},
				},
				{
					FileID:    "file-2.mp4",
					Size:      2222,
					Name:      "file-2.mp4",
					CreatedAt: time.Time{},
				},
				{
					FileID:    "file-3.mp4",
					Size:      3333,
					Name:      "file-3.mp4",
					CreatedAt: time.Time{},
				},
			},
			wantErr: false,
		},
		{
			name: "failed to get all files",
			args: args{
				ctx: context.Background(),
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().GetAllFiles(context.Background()).Return([]dbstore.FileDetail{}, fmt.Errorf("some-err"))
			},
			want:    []FileInfo{},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDBStore := dbStoreMocks.NewMockDBStore(ctrl)

			tt.mockFunc(mockDBStore)

			s := service{
				dbStore: mockDBStore,
			}
			got, err := s.GetAllFiles(tt.args.ctx)
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

func Test_service_DeleteFileByID(t *testing.T) {
	type args struct {
		ctx context.Context
		id  string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockDBStore *dbStoreMocks.MockDBStore)
		wantErr  bool
	}{
		{
			name: "successfully delete a file by its ID",
			args: args{
				ctx: context.Background(),
				id:  "some-id",
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				os.Create("storage/videos/file-id")
				mockDBStore.EXPECT().DeleteFileByID(context.Background(), "some-id").Return(dbstore.FileDetail{
					ID:        "file-id",
					Size:      123,
					Path:      "path/to/file-id",
					CreatedAt: time.Time{},
				}, nil)
			},
			wantErr: false,
		},
		{
			name: "failed to delete from DB",
			args: args{
				ctx: context.Background(),
				id:  "some-id",
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().DeleteFileByID(context.Background(), "some-id").Return(dbstore.FileDetail{}, fmt.Errorf("some-err"))
			},
			wantErr: true,
		},
		{
			name: "failed to delete file from system",
			args: args{
				ctx: context.Background(),
				id:  "some-id",
			},
			mockFunc: func(mockDBStore *dbStoreMocks.MockDBStore) {
				mockDBStore.EXPECT().DeleteFileByID(context.Background(), "some-id").Return(dbstore.FileDetail{
					ID:        "file-id",
					Size:      123,
					Path:      "path/to/file-id",
					CreatedAt: time.Time{},
				}, nil)
				mockDBStore.EXPECT().InsertNewFile(context.Background(), dbstore.FileDetail{
					ID:        "file-id",
					Size:      123,
					Path:      "path/to/file-id",
					CreatedAt: time.Time{},
				}).Return(nil)
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockDBStore := dbStoreMocks.NewMockDBStore(ctrl)

			tt.mockFunc(mockDBStore)

			s := service{
				dbStore: mockDBStore,
			}
			if err := s.DeleteFileByID(tt.args.ctx, tt.args.id); (err != nil) != tt.wantErr {
				t.Errorf("DeleteFileByID() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
