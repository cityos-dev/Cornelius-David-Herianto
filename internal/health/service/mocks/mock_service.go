// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cityos-dev/Cornelius-David-Herianto/internal/health/service (interfaces: Service)

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
)

// MockService is a mock of Service interface.
type MockService struct {
	ctrl     *gomock.Controller
	recorder *MockServiceMockRecorder
}

// MockServiceMockRecorder is the mock recorder for MockService.
type MockServiceMockRecorder struct {
	mock *MockService
}

// NewMockService creates a new mock instance.
func NewMockService(ctrl *gomock.Controller) *MockService {
	mock := &MockService{ctrl: ctrl}
	mock.recorder = &MockServiceMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockService) EXPECT() *MockServiceMockRecorder {
	return m.recorder
}

// GetServiceHealth mocks base method.
func (m *MockService) GetServiceHealth(arg0 context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetServiceHealth", arg0)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetServiceHealth indicates an expected call of GetServiceHealth.
func (mr *MockServiceMockRecorder) GetServiceHealth(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetServiceHealth", reflect.TypeOf((*MockService)(nil).GetServiceHealth), arg0)
}
