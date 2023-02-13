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

type Service interface {
	UploadFile(ctx context.Context, file multipart.File, host, filename string, size int64) (string, error)
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
