package handler

import (
	"fmt"
	"net/http"

	"github.com/labstack/echo"

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
	defer multipartFile.Close()

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

	ctx.Response().Header().Set("Content-Disposition", fmt.Sprintf("form-data; name='data'; filename=%s", fileID))
	return ctx.File("storage/videos/" + fileID)
}

func (h filesHTTPHandler) GetAllFiles(ctx echo.Context) error {
	files, err := h.service.GetAllFiles(ctx.Request().Context())
	if err != nil {
		return echo.NewHTTPError(http.StatusInternalServerError, "failed to get all files from DB", err)
	}
	return ctx.JSON(http.StatusOK, files)
}
