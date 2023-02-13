package service

import (
	"context"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"time"

	"github.com/lib/pq"
	"golang.org/x/exp/slices"

	"github.com/cityos-dev/Cornelius-David-Herianto/helper/uuid"
	filesDBStore "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore"
)

const (
	localStoragePath = "storage/videos/"
)

var (
	ErrorUnsupportedFileTypes = fmt.Errorf("unsupported file types")
	ErrorDuplicateKey         = fmt.Errorf("duplicate key value")
)

var allowedContentType = []string{
	"video/mp4",
	"video/mpeg",
}

type FileInfo struct {
	FileID    string    `json:"fileid"`
	Name      string    `json:"name"`
	Size      int64     `json:"size"`
	CreatedAt time.Time `json:"created_at"`
}

type Service interface {
	UploadFile(ctx context.Context, file multipart.File, host, filename string, size int64) (string, error)
	GetAllFiles(ctx context.Context) ([]FileInfo, error)
}

type service struct {
	dbStore   filesDBStore.DBStore
	uuidUtils uuid.Utils
}

func New(dbStore filesDBStore.DBStore, uuidUtils uuid.Utils) Service {
	return service{
		dbStore:   dbStore,
		uuidUtils: uuidUtils,
	}
}

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
	defer dst.Close()

	fileBytes, _ := io.ReadAll(file)

	// get content type from header
	contentType := http.DetectContentType(fileBytes[:512])

	if !slices.Contains(allowedContentType, contentType) {
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

func (s service) GetAllFiles(ctx context.Context) ([]FileInfo, error) {
	files, err := s.dbStore.GetAllFiles(ctx)
	if err != nil {
		return []FileInfo{}, fmt.Errorf("failed to get all files from DB, err: %v", err)
	}
	var fileInfos []FileInfo
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