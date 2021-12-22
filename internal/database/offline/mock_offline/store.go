// Code generated by MockGen. DO NOT EDIT.
// Source: internal/database/offline/store.go

// Package mock_offline is a generated GoMock package.
package mock_offline

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	offline "github.com/oom-ai/oomstore/internal/database/offline"
	types "github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// MockStore is a mock of Store interface.
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore.
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance.
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Close mocks base method.
func (m *MockStore) Close() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Close")
	ret0, _ := ret[0].(error)
	return ret0
}

// Close indicates an expected call of Close.
func (mr *MockStoreMockRecorder) Close() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Close", reflect.TypeOf((*MockStore)(nil).Close))
}

// Export mocks base method.
func (m *MockStore) Export(ctx context.Context, opt offline.ExportOpt) (<-chan types.ExportRecord, <-chan error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Export", ctx, opt)
	ret0, _ := ret[0].(<-chan types.ExportRecord)
	ret1, _ := ret[1].(<-chan error)
	return ret0, ret1
}

// Export indicates an expected call of Export.
func (mr *MockStoreMockRecorder) Export(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Export", reflect.TypeOf((*MockStore)(nil).Export), ctx, opt)
}

// Import mocks base method.
func (m *MockStore) Import(ctx context.Context, opt offline.ImportOpt) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Import", ctx, opt)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Import indicates an expected call of Import.
func (mr *MockStoreMockRecorder) Import(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Import", reflect.TypeOf((*MockStore)(nil).Import), ctx, opt)
}

// Join mocks base method.
func (m *MockStore) Join(ctx context.Context, opt offline.JoinOpt) (*types.JoinResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Join", ctx, opt)
	ret0, _ := ret[0].(*types.JoinResult)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Join indicates an expected call of Join.
func (mr *MockStoreMockRecorder) Join(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Join", reflect.TypeOf((*MockStore)(nil).Join), ctx, opt)
}

// Ping mocks base method.
func (m *MockStore) Ping(ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Ping", ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// Ping indicates an expected call of Ping.
func (mr *MockStoreMockRecorder) Ping(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Ping", reflect.TypeOf((*MockStore)(nil).Ping), ctx)
}

// TableSchema mocks base method.
func (m *MockStore) TableSchema(ctx context.Context, tableName string) (*types.DataTableSchema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableSchema", ctx, tableName)
	ret0, _ := ret[0].(*types.DataTableSchema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TableSchema indicates an expected call of TableSchema.
func (mr *MockStoreMockRecorder) TableSchema(ctx, tableName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableSchema", reflect.TypeOf((*MockStore)(nil).TableSchema), ctx, tableName)
}
