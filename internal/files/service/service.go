package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"os"
	"path/filepath"
	"time"

	"github.com/lib/pq"
	"golang.org/x/exp/slices"

	filesDBStore "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore"
)

const (
	localStoragePath = "storage/videos/"
)

// Errors represent custom error that will be verified by the handler layer
var (
	ErrorUnsupportedFileTypes = fmt.Errorf("unsupported file types")
	ErrorDuplicateKey         = fmt.Errorf("duplicate key value")
)

var allowedExtensions = []string{
	".mp4", ".mpg", ".mpeg",
}

// FileInfo represents information of a file
type FileInfo struct {
	FileID    string    `json:"fileid"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

// Service provides mechanism to interact with files
//
//go:generate mockgen -destination mocks/mock_service.go github.com/cityos-dev/Cornelius-David-Herianto/internal/files/service Service
type Service interface {
	UploadFile(ctx context.Context, file multipart.File, host, filename string, size int64) (string, error)
	GetAllFiles(ctx context.Context) ([]FileInfo, error)
	DeleteFileByID(ctx context.Context, id string) error
}

type service struct {
	dbStore filesDBStore.DBStore
}

// New returned new Service instance
func New(dbStore filesDBStore.DBStore) Service {
	return service{
		dbStore: dbStore,
	}
}

// UploadFile do save file to local storage (file system) and also insert the file detail info to the DB
func (s service) UploadFile(ctx context.Context, file multipart.File, host, filename string, size int64) (string, error) {
	// save file to local storage
	err := os.MkdirAll(localStoragePath, os.ModePerm)
	if err != nil {
		return "", fmt.Errorf("failed to create directory: %s, err: %v", localStoragePath, err)
	}
	targetFilename := localStoragePath + filename
	dst, err := os.Create(targetFilename)
	if err != nil {
		return "", fmt.Errorf("failed to create file: %s, err: %v", targetFilename, err)
	}

	defer func() {
		_ = dst.Close()
	}()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return "", fmt.Errorf("failed to write read uploaded file, err: %v", err)
	}

	// validate content type
	if !slices.Contains(allowedExtensions, filepath.Ext(filename)) {
		return "", ErrorUnsupportedFileTypes
	}

	_, err = dst.Write(fileBytes)
	if err != nil {
		return "", fmt.Errorf("failed to write file to local storage, err: %v", err)
	}

	fileFullPath := host + "/v1/files/" + filename

	err = s.dbStore.InsertNewFile(ctx, filesDBStore.FileDetail{
		ID:        filename,
		Size:      size,
		Path:      fileFullPath,
		CreatedAt: time.Time{},
	})
	if err != nil {
		if pgErr, ok := err.(*pq.Error); ok {
			if pgErr.Code == "23505" {
				return "", ErrorDuplicateKey
			}
		}
		// revert file saving
		_ = os.Remove(localStoragePath + filename)
		return "", fmt.Errorf("failed to insert file information to DB, err: %v", err)
	}

	return fileFullPath, nil
}

// GetAllFiles returned all files info listed on the DB
func (s service) GetAllFiles(ctx context.Context) ([]FileInfo, error) {
	files, err := s.dbStore.GetAllFiles(ctx)
	if err != nil {
		return []FileInfo{}, fmt.Errorf("failed to get all files from DB, err: %v", err)
	}
	fileInfos := make([]FileInfo, 0)
	for _, file := range files {
		fileInfos = append(fileInfos, mapFileDetailsToFileInfo(file))
	}
	return fileInfos, nil
}

func mapFileDetailsToFileInfo(fileDetail filesDBStore.FileDetail) FileInfo {
	return FileInfo{
		FileID:    fileDetail.ID,
		Name:      fileDetail.ID,
		Size:      fileDetail.Size,
		CreatedAt: fileDetail.CreatedAt,
	}
}

// DeleteFileByID delete a file by its id
func (s service) DeleteFileByID(ctx context.Context, id string) error {
	fileDetail, err := s.dbStore.DeleteFileByID(ctx, id)
	if err != nil {
		return fmt.Errorf("failed to delete entry from DB, err: %w", err)
	}

	err = os.Remove(localStoragePath + fileDetail.ID)
	if err != nil {
		_ = s.dbStore.InsertNewFile(ctx, fileDetail)
		return fmt.Errorf("failed to delete file from local storage, err: %w", err)
	}
	return nil
}
