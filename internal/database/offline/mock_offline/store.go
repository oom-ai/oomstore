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

// CreateTable mocks base method.
func (m *MockStore) CreateTable(ctx context.Context, opt offline.CreateTableOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateTable", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateTable indicates an expected call of CreateTable.
func (mr *MockStoreMockRecorder) CreateTable(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateTable", reflect.TypeOf((*MockStore)(nil).CreateTable), ctx, opt)
}

// Export mocks base method.
func (m *MockStore) Export(ctx context.Context, opt offline.ExportOpt) (*types.ExportResult, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Export", ctx, opt)
	ret0, _ := ret[0].(*types.ExportResult)
	ret1, _ := ret[1].(error)
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

// Push mocks base method.
func (m *MockStore) Push(ctx context.Context, opt offline.PushOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Push", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Push indicates an expected call of Push.
func (mr *MockStoreMockRecorder) Push(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Push", reflect.TypeOf((*MockStore)(nil).Push), ctx, opt)
}

// Snapshot mocks base method.
func (m *MockStore) Snapshot(ctx context.Context, opt offline.SnapshotOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Snapshot", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// Snapshot indicates an expected call of Snapshot.
func (mr *MockStoreMockRecorder) Snapshot(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Snapshot", reflect.TypeOf((*MockStore)(nil).Snapshot), ctx, opt)
}

// TableSchema mocks base method.
func (m *MockStore) TableSchema(ctx context.Context, opt offline.TableSchemaOpt) (*types.DataTableSchema, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "TableSchema", ctx, opt)
	ret0, _ := ret[0].(*types.DataTableSchema)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// TableSchema indicates an expected call of TableSchema.
func (mr *MockStoreMockRecorder) TableSchema(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "TableSchema", reflect.TypeOf((*MockStore)(nil).TableSchema), ctx, opt)
}
