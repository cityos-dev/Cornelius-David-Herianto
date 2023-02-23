// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore (interfaces: DBStore)

// Package mock_dbstore is a generated GoMock package.
package mock_dbstore

import (
	context "context"
	reflect "reflect"

	dbstore "github.com/cityos-dev/Cornelius-David-Herianto/internal/files/store/dbstore"
	gomock "github.com/golang/mock/gomock"
)

// MockDBStore is a mock of DBStore interface.
type MockDBStore struct {
	ctrl     *gomock.Controller
	recorder *MockDBStoreMockRecorder
}

// MockDBStoreMockRecorder is the mock recorder for MockDBStore.
type MockDBStoreMockRecorder struct {
	mock *MockDBStore
}

// NewMockDBStore creates a new mock instance.
func NewMockDBStore(ctrl *gomock.Controller) *MockDBStore {
	mock := &MockDBStore{ctrl: ctrl}
	mock.recorder = &MockDBStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockDBStore) EXPECT() *MockDBStoreMockRecorder {
	return m.recorder
}

// DeleteFileByID mocks base method.
func (m *MockDBStore) DeleteFileByID(arg0 context.Context, arg1 string) (dbstore.FileDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteFileByID", arg0, arg1)
	ret0, _ := ret[0].(dbstore.FileDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DeleteFileByID indicates an expected call of DeleteFileByID.
func (mr *MockDBStoreMockRecorder) DeleteFileByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteFileByID", reflect.TypeOf((*MockDBStore)(nil).DeleteFileByID), arg0, arg1)
}

// GetAllFiles mocks base method.
func (m *MockDBStore) GetAllFiles(arg0 context.Context) ([]dbstore.FileDetail, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllFiles", arg0)
	ret0, _ := ret[0].([]dbstore.FileDetail)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllFiles indicates an expected call of GetAllFiles.
func (mr *MockDBStoreMockRecorder) GetAllFiles(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllFiles", reflect.TypeOf((*MockDBStore)(nil).GetAllFiles), arg0)
}

// InsertNewFile mocks base method.
func (m *MockDBStore) InsertNewFile(arg0 context.Context, arg1 dbstore.FileDetail) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "InsertNewFile", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// InsertNewFile indicates an expected call of InsertNewFile.
func (mr *MockDBStoreMockRecorder) InsertNewFile(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "InsertNewFile", reflect.TypeOf((*MockDBStore)(nil).InsertNewFile), arg0, arg1)
}