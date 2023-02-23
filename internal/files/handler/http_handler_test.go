package handler

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"reflect"
	"strings"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	"github.com/cityos-dev/Cornelius-David-Herianto/internal/files/service"
	filesSvc "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/service"
	filesSvcMock "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/service/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		service service.Service
	}
	tests := []struct {
		name string
		args args
		want filesHTTPHandler
	}{
		{
			name: "successfully get new files http handler",
			args: args{
				service: nil,
			},
			want: filesHTTPHandler{
				service: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.service); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_filesHTTPHandler_DeleteFileByID(t *testing.T) {
	type args struct {
		method string
		url    string
	}
	type want struct {
		body        string
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockService *filesSvcMock.MockService)
		want     want
		wantErr  bool
	}{
		{
			name: "successfully delete a file by its id",
			args: args{
				method: http.MethodDelete,
				url:    "http://localhost/v1/files/test.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().DeleteFileByID(gomock.Any(), "test.mp4").Return(nil)
			},
			want: want{
				body:        `OK`,
				code:        http.StatusNoContent,
				contentType: "text/plain; charset=UTF-8",
			},
			wantErr: false,
		},
		{
			name: "deleted id not found",
			args: args{
				method: http.MethodDelete,
				url:    "http://localhost/v1/files/test.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().DeleteFileByID(gomock.Any(), "test.mp4").Return(sql.ErrNoRows)
			},
			want: want{
				body: `{"message":"deleted file is not exists","dev_message":"sql: no rows in result set"}`,
				code: http.StatusNotFound,
			},
			wantErr: true,
		},
		{
			name: "failed to delete a file by its id (other error)",
			args: args{
				method: http.MethodDelete,
				url:    "http://localhost/v1/files/test.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().DeleteFileByID(gomock.Any(), "test.mp4").Return(fmt.Errorf("some-err"))
			},
			want: want{
				body: `{"message":"failed to delete file with id: test.mp4","dev_message":"some-err"}`,
				code: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockFilesSvc := filesSvcMock.NewMockService(ctrl)
			tt.mockFunc(mockFilesSvc)

			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			ctx := echo.New().NewContext(r, w)
			ctx.SetPath("v1/files/:fileID")
			ctx.SetParamNames("fileID")
			urlPaths := strings.Split(tt.args.url, "/")
			ctx.SetParamValues(urlPaths[len(urlPaths)-1])

			h := filesHTTPHandler{
				service: mockFilesSvc,
			}

			err := h.DeleteFileByID(ctx)
			if tt.wantErr {
				httpErr := err.(*echo.HTTPError)
				if httpErr.Code != tt.want.code {
					t.Errorf("DeleteFileByID() status code got = %d, want %d\n", httpErr.Code, tt.want.code)
				}
				errMsgByte, _ := json.Marshal(httpErr.Message)
				if strings.TrimSpace(string(errMsgByte)) != tt.want.body {
					t.Errorf("DeleteFileByID() body got = %s, want %s\n", string(errMsgByte), tt.want.body)
				}
				return
			}

			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			if res.StatusCode != tt.want.code {
				t.Errorf("DeleteFileByID() status code got = %d, want %d\n", res.StatusCode, tt.want.code)
			}

			if res.Header.Get(echo.HeaderContentType) != tt.want.contentType {
				t.Errorf("DeleteFileByID() content-type got = %s, want %s\n", res.Header.Get(echo.HeaderContentType), tt.want.contentType)
			}

			if err != nil {
				t.Errorf("WriteResponse DeleteFileByID() read from body err = %v\n", err)
			}

			if strings.TrimSpace(string(resBody)) != tt.want.body {
				t.Errorf("DeleteFileByID() body got = %s, want %s\n", string(resBody), tt.want.body)
			}
		})
	}
}

func Test_filesHTTPHandler_GetAllFiles(t *testing.T) {
	type args struct {
		method string
		url    string
	}
	type want struct {
		body        string
		code        int
		contentType string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockService *filesSvcMock.MockService)
		want     want
		wantErr  bool
	}{
		{
			name: "successfully get all files but no result",
			args: args{
				method: http.MethodGet,
				url:    "http://localhost/v1/files",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().GetAllFiles(gomock.Any()).Return([]filesSvc.FileInfo{}, nil)
			},
			want: want{
				body:        "[]",
				code:        http.StatusOK,
				contentType: "application/json; charset=UTF-8",
			},
			wantErr: false,
		},
		{
			name: "successfully get all files (got 3 files)",
			args: args{
				method: http.MethodGet,
				url:    "http://localhost/v1/files",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().GetAllFiles(gomock.Any()).Return([]filesSvc.FileInfo{
					{
						FileID:    "file-1.mp4",
						Name:      "file-1.mp4",
						Size:      111,
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						FileID:    "file-2.mp4",
						Name:      "file-2.mp4",
						Size:      222,
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
					{
						FileID:    "file-3.mp4",
						Name:      "file-3.mp4",
						Size:      333,
						CreatedAt: time.Date(2023, 1, 1, 0, 0, 0, 0, time.UTC),
					},
				}, nil)
			},
			want: want{
				body:        `[{"fileid":"file-1.mp4","name":"file-1.mp4","size":111,"created_at":"2023-01-01T00:00:00Z"},{"fileid":"file-2.mp4","name":"file-2.mp4","size":222,"created_at":"2023-01-01T00:00:00Z"},{"fileid":"file-3.mp4","name":"file-3.mp4","size":333,"created_at":"2023-01-01T00:00:00Z"}]`,
				code:        http.StatusOK,
				contentType: "application/json; charset=UTF-8",
			},
			wantErr: false,
		},
		{
			name: "failed to get all files",
			args: args{
				method: http.MethodGet,
				url:    "http://localhost/v1/files",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().GetAllFiles(gomock.Any()).Return([]filesSvc.FileInfo{}, fmt.Errorf("some-err"))
			},
			want: want{
				body:        `{"message":"failed to get all files from DB","dev_message":"some-err"}`,
				code:        http.StatusInternalServerError,
				contentType: "application/json; charset=UTF-8",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockFilesSvc := filesSvcMock.NewMockService(ctrl)
			tt.mockFunc(mockFilesSvc)

			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			ctx := echo.New().NewContext(r, w)

			h := filesHTTPHandler{
				service: mockFilesSvc,
			}

			err := h.GetAllFiles(ctx)
			if tt.wantErr {
				httpErr := err.(*echo.HTTPError)
				if httpErr.Code != tt.want.code {
					t.Errorf("GetAllFiles() status code got = %d, want %d\n", httpErr.Code, tt.want.code)
				}
				errMsgByte, _ := json.Marshal(httpErr.Message)
				if strings.TrimSpace(string(errMsgByte)) != tt.want.body {
					t.Errorf("GetAllFiles() body got = %s, want %s\n", string(errMsgByte), tt.want.body)
				}
				return
			}

			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			if res.StatusCode != tt.want.code {
				t.Errorf("GetAllFiles() status code got = %d, want %d\n", res.StatusCode, tt.want.code)
			}

			if res.Header.Get(echo.HeaderContentType) != tt.want.contentType {
				t.Errorf("GetAllFiles() content-type got = %s, want %s\n", res.Header.Get(echo.HeaderContentType), tt.want.contentType)
			}

			if err != nil {
				t.Errorf("WriteResponse GetAllFiles() read from body err = %v\n", err)
			}

			if strings.TrimSpace(string(resBody)) != tt.want.body {
				t.Errorf("GetAllFiles() body got = %s, want %s\n", string(resBody), tt.want.body)
			}
		})
	}
}

func Test_filesHTTPHandler_GetFileByID(t *testing.T) {
	type args struct {
		method string
		url    string
	}
	type want struct {
		body               string
		code               int
		contentType        string
		contentDisposition string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockService *filesSvcMock.MockService)
		want     want
		wantErr  bool
	}{
		{
			name: "successfully get the requested file",
			args: args{
				method: http.MethodGet,
				url:    "http://localhost/v1/files/sample.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {

			},
			want: want{
				code:               http.StatusOK,
				contentType:        "video/mp4",
				contentDisposition: "form-data; name='data'; filename=sample.mp4",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockFilesSvc := filesSvcMock.NewMockService(ctrl)
			tt.mockFunc(mockFilesSvc)

			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			ctx := echo.New().NewContext(r, w)
			ctx.SetPath("v1/files/:fileID")
			ctx.SetParamNames("fileID")
			urlPaths := strings.Split(tt.args.url, "/")
			ctx.SetParamValues(urlPaths[len(urlPaths)-1])

			h := filesHTTPHandler{
				service: mockFilesSvc,
			}

			err := h.GetFileByID(ctx)
			if tt.wantErr {
				httpErr := err.(*echo.HTTPError)
				if httpErr.Code != tt.want.code {
					t.Errorf("GetFileByID() status code got = %d, want %d\n", httpErr.Code, tt.want.code)
				}
				errMsgByte, _ := json.Marshal(httpErr.Message)
				if strings.TrimSpace(string(errMsgByte)) != tt.want.body {
					t.Errorf("GetFileByID() body got = %s, want %s\n", string(errMsgByte), tt.want.body)
				}
				return
			}

			res := w.Result()
			defer res.Body.Close()

			if res.StatusCode != tt.want.code {
				t.Errorf("GetFileByID() status code got = %d, want %d\n", res.StatusCode, tt.want.code)
			}

			if res.Header.Get(echo.HeaderContentType) != tt.want.contentType {
				t.Errorf("GetFileByID() content-type got = %s, want %s\n", res.Header.Get(echo.HeaderContentType), tt.want.contentType)
			}
			if res.Header.Get(echo.HeaderContentDisposition) != tt.want.contentDisposition {
				t.Errorf("GetFileByID() content-disposition got = %s, want %s\n", res.Header.Get(echo.HeaderContentDisposition), tt.want.contentDisposition)
			}
		})
	}
}

func Test_filesHTTPHandler_UploadFile(t *testing.T) {
	type args struct {
		method   string
		url      string
		filepath string
	}
	type want struct {
		body        string
		code        int
		contentType string
		location    string
	}
	tests := []struct {
		name     string
		args     args
		mockFunc func(mockService *filesSvcMock.MockService)
		want     want
		wantErr  bool
	}{
		{
			name: "successfully uploaded a file",
			args: args{
				method:   http.MethodPost,
				url:      "http://localhost/v1/files",
				filepath: "test/post_1/sample.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), "localhost", "sample.mp4", int64(2848208)).Return("localhost/v1/files/sample.mpg", nil)
			},
			want: want{
				body:        `OK`,
				code:        http.StatusCreated,
				contentType: "text/plain; charset=UTF-8",
				location:    "localhost/v1/files/sample.mpg",
			},
			wantErr: false,
		},
		{
			name: "upload unsupported file type",
			args: args{
				method:   http.MethodPost,
				url:      "http://localhost/v1/files",
				filepath: "test/post_4/test.txt",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), "localhost", "test.txt", int64(19)).Return("", filesSvc.ErrorUnsupportedFileTypes)
			},
			want: want{
				body: `{"message":"invalid content type, only video/mp4 and video/mpeg allowed","dev_message":"unsupported file types"}`,
				code: http.StatusUnsupportedMediaType,
			},
			wantErr: true,
		},
		{
			name: "duplicated key",
			args: args{
				method:   http.MethodPost,
				url:      "http://localhost/v1/files",
				filepath: "test/post_1/sample.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), "localhost", "sample.mp4", int64(2848208)).Return("", filesSvc.ErrorDuplicateKey)
			},
			want: want{
				body: `{"message":"file with id: sample.mp4 is already exist","dev_message":"duplicate key value"}`,
				code: http.StatusConflict,
			},
			wantErr: true,
		},
		{
			name: "failed to upload file (other error)",
			args: args{
				method:   http.MethodPost,
				url:      "http://localhost/v1/files",
				filepath: "test/post_1/sample.mp4",
			},
			mockFunc: func(mockService *filesSvcMock.MockService) {
				mockService.EXPECT().UploadFile(gomock.Any(), gomock.Any(), "localhost", "sample.mp4", int64(2848208)).Return("", fmt.Errorf("some-err"))
			},
			want: want{
				body: `{"message":"failed to upload the file, please try again later","dev_message":"some-err"}`,
				code: http.StatusInternalServerError,
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockFilesSvc := filesSvcMock.NewMockService(ctrl)

			tt.mockFunc(mockFilesSvc)

			filePath := "../../../" + tt.args.filepath
			body := new(bytes.Buffer)
			mw := multipart.NewWriter(body)
			file, err := os.Open(filePath)
			if err != nil {
				t.Fatal(err)
			}
			writer, err := mw.CreateFormFile("data", filePath)
			if err != nil {
				t.Fatal(err)
			}
			if _, err := io.Copy(writer, file); err != nil {
				t.Fatal(err)
			}
			_ = mw.Close()

			r := httptest.NewRequest(tt.args.method, tt.args.url, body)
			r.Header.Add(echo.HeaderContentType, mw.FormDataContentType())
			w := httptest.NewRecorder()
			ctx := echo.New().NewContext(r, w)

			h := filesHTTPHandler{
				service: mockFilesSvc,
			}

			err = h.UploadFile(ctx)
			if tt.wantErr {
				httpErr := err.(*echo.HTTPError)
				if httpErr.Code != tt.want.code {
					t.Errorf("UploadFile() status code got = %d, want %d\n", httpErr.Code, tt.want.code)
				}
				errMsgByte, _ := json.Marshal(httpErr.Message)
				if strings.TrimSpace(string(errMsgByte)) != tt.want.body {
					t.Errorf("UploadFile() body got = %s, want %s\n", string(errMsgByte), tt.want.body)
				}
				return
			}

			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			if res.StatusCode != tt.want.code {
				t.Errorf("UploadFile() status code got = %d, want %d\n", res.StatusCode, tt.want.code)
			}

			if res.Header.Get(echo.HeaderContentType) != tt.want.contentType {
				t.Errorf("UploadFile() content-type got = %s, want %s\n", res.Header.Get(echo.HeaderContentType), tt.want.contentType)
			}

			if res.Header.Get("Location") != tt.want.location {
				t.Errorf("UploadFile() location got = %s, want %s\n", res.Header.Get("Location"), tt.want.location)
			}

			if err != nil {
				t.Errorf("WriteResponse UploadFile() read from body err = %v\n", err)
			}

			if strings.TrimSpace(string(resBody)) != tt.want.body {
				t.Errorf("UploadFile() body got = %s, want %s\n", string(resBody), tt.want.body)
			}
		})
	}
}
