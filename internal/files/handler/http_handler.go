package handler

import (
	"database/sql"
	"errors"
	"fmt"
	"mime"
	"net/http"
	"path/filepath"

	"github.com/labstack/echo/v4"

	httpHelper "github.com/cityos-dev/Cornelius-David-Herianto/helper/http"
	filesSvc "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/service"
)

type filesHTTPHandler struct {
	service filesSvc.Service
}

func New(service filesSvc.Service) filesHTTPHandler {
	return filesHTTPHandler{
		service: service,
	}
}

func (h filesHTTPHandler) UploadFile(ctx echo.Context) error {
	multipartFileHeader, err := ctx.FormFile("data")
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httpHelper.NewErrorMessage("failed to process uploaded file", err))
	}
	multipartFile, err := multipartFileHeader.Open()
	if err != nil {
		return echo.NewHTTPError(http.StatusBadRequest, httpHelper.NewErrorMessage("failed to process uploaded file", err))
	}
	defer func() {
		_ = multipartFile.Close()
	}()

	location, err := h.service.UploadFile(ctx.Request().Context(), multipartFile, ctx.Request().Host, multipartFileHeader.Filename, multipartFileHeader.Size)
	if err != nil {
		if err == filesSvc.ErrorUnsupportedFileTypes {
			return echo.NewHTTPError(http.StatusUnsupportedMediaType, httpHelper.NewErrorMessage("invalid content type, only video/mp4 and video/mpeg allowed", err))
		} else if err == filesSvc.ErrorDuplicateKey {
			return echo.NewHTTPError(http.StatusConflict, httpHelper.NewErrorMessage(fmt.Sprintf("file with id: %s is already exist", multipartFileHeader.Filename), err))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, httpHelper.NewErrorMessage("failed to upload the file, please try again later", err))
	}

	ctx.Response().Header().Set("Location", location)
	return ctx.String(http.StatusCreated, "OK")
}

func (h filesHTTPHandler) GetFileByID(ctx echo.Context) error {
	fileID := ctx.Param("fileID")

	ctx.Response().Header().Set(echo.HeaderContentDisposition, fmt.Sprintf("form-data; name='data'; filename=%s", fileID))
	ctx.Response().Header().Set(echo.HeaderContentType, mime.TypeByExtension(filepath.Ext(fileID)))
	return ctx.File("storage/videos/" + fileID)
}

func (h filesHTTPHandler) GetAllFiles(ctx echo.Context) error {
	files, err := h.service.GetAllFiles(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, httpHelper.NewErrorMessage("failed to get all files from DB", err))
	}
	return ctx.JSON(http.StatusOK, files)
}

func (h filesHTTPHandler) DeleteFileByID(ctx echo.Context) error {
	err := h.service.DeleteFileByID(ctx.Request().Context(), ctx.Param("fileID"))
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return echo.NewHTTPError(http.StatusNotFound, httpHelper.NewErrorMessage("deleted file is not exists", err))
		}
		return echo.NewHTTPError(http.StatusInternalServerError, httpHelper.NewErrorMessage(fmt.Sprintf("failed to delete file with id: %s", ctx.Param("fileID")), err))
	}
	return ctx.String(http.StatusNoContent, "OK")
}
