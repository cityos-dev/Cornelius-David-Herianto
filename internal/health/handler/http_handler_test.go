package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"reflect"
	"strings"
	"testing"

	"github.com/golang/mock/gomock"
	"github.com/labstack/echo/v4"

	"github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service"
	healthSvcMock "github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service/mocks"
)

func TestNew(t *testing.T) {
	type args struct {
		healthSvc service.Service
	}
	tests := []struct {
		name string
		args args
		want healthHTTPHandler
	}{
		{
			name: "success get new httpHandler instance",
			args: args{
				healthSvc: nil,
			},
			want: healthHTTPHandler{
				healthSvc: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.healthSvc); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_healthHTTPHandler_GetHealth(t *testing.T) {
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
		mockFunc func(mockService *healthSvcMock.MockService)
		want     want
		wantErr  bool
	}{
		{
			name: "service health ok",
			args: args{
				method: http.MethodGet,
				url:    "http://localhost/v1/health",
			},
			mockFunc: func(mockService *healthSvcMock.MockService) {
				mockService.EXPECT().GetServiceHealth(context.Background()).Return(nil)
			},
			want: want{
				body:        `OK`,
				code:        http.StatusOK,
				contentType: "text/plain; charset=UTF-8",
			},
			wantErr: false,
		},
		{
			name: "service health not ok",
			args: args{
				method: http.MethodGet,
				url:    "http://localhost/v1/health",
			},
			mockFunc: func(mockService *healthSvcMock.MockService) {
				mockService.EXPECT().GetServiceHealth(context.Background()).Return(fmt.Errorf("some-error"))
			},
			want: want{
				body:        `{"message":"service is not healthy","dev_message":"some-error"}`,
				code:        http.StatusInternalServerError,
				contentType: "application/json",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			mockHealthSvc := healthSvcMock.NewMockService(ctrl)

			tt.mockFunc(mockHealthSvc)

			r := httptest.NewRequest(tt.args.method, tt.args.url, nil)
			w := httptest.NewRecorder()
			ctx := echo.New().NewContext(r, w)

			h := &healthHTTPHandler{
				healthSvc: mockHealthSvc,
			}

			err := h.GetHealth(ctx)
			if tt.wantErr {
				httpErr := err.(*echo.HTTPError)
				if httpErr.Code != tt.want.code {
					t.Errorf("GetHealth() status code got = %d, want %d\n", httpErr.Code, tt.want.code)
				}
				errMsgByte, _ := json.Marshal(httpErr.Message)
				if strings.TrimSpace(string(errMsgByte)) != tt.want.body {
					t.Errorf("GetHealth() body got = %s, want %s\n", string(errMsgByte), tt.want.body)
				}
				return
			}

			res := w.Result()
			defer res.Body.Close()
			resBody, err := io.ReadAll(res.Body)

			if res.StatusCode != tt.want.code {
				t.Errorf("GetHealth() status code got = %d, want %d\n", res.StatusCode, tt.want.code)
			}

			if res.Header.Get(echo.HeaderContentType) != tt.want.contentType {
				t.Errorf("GetHealth() content-type got = %s, want %s\n", res.Header.Get(echo.HeaderContentType), tt.want.contentType)
			}

			if err != nil {
				t.Errorf("WriteResponse GetHealth() read from body err = %v\n", err)
			}

			if strings.TrimSpace(string(resBody)) != tt.want.body {
				t.Errorf("GetHealth() body got = %s, want %s\n", string(resBody), tt.want.body)
			}
		})
	}
}
