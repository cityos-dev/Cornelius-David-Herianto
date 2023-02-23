package service

import (
	"context"
	"fmt"
	"reflect"
	"testing"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/jmoiron/sqlx"
)

func Test_service_GetServiceHealth(t *testing.T) {
	tests := []struct {
		name     string
		mockFunc func(sqlMock sqlmock.Sqlmock)
		wantErr  bool
	}{
		{
			name: "all goes well",
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectPing().WillReturnError(nil)
			},
			wantErr: false,
		},
		{
			name: "failed get health",
			mockFunc: func(sqlMock sqlmock.Sqlmock) {
				sqlMock.ExpectPing().WillReturnError(fmt.Errorf("some-error"))
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			mockDB, sqlMock, err := sqlmock.New(sqlmock.MonitorPingsOption(true))
			if err != nil {
				t.Errorf("error when opening a database connection: %v\n", err)
			}
			defer mockDB.Close()
			tt.mockFunc(sqlMock)
			s := &service{
				dbConn: sqlx.NewDb(mockDB, "postgres"),
			}
			if err := s.GetServiceHealth(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("GetServiceHealth() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func TestNew(t *testing.T) {
	type args struct {
		dbConn *sqlx.DB
	}
	tests := []struct {
		name string
		args args
		want Service
	}{
		{
			name: "success get new Service instance",
			args: args{
				dbConn: nil,
			},
			want: &service{
				dbConn: nil,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.dbConn); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("New() = %v, want %v", got, tt.want)
			}
		})
	}
}
