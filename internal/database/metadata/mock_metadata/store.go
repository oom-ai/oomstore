// Code generated by MockGen. DO NOT EDIT.
// Source: internal/database/metadata/store.go

// Package mock_metadata is a generated GoMock package.
package mock_metadata

import (
	context "context"
	sql "database/sql"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	metadata "github.com/oom-ai/oomstore/internal/database/metadata"
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

// CreateEntity mocks base method.
func (m *MockStore) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntity", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntity indicates an expected call of CreateEntity.
func (mr *MockStoreMockRecorder) CreateEntity(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntity", reflect.TypeOf((*MockStore)(nil).CreateEntity), ctx, opt)
}

// CreateFeature mocks base method.
func (m *MockStore) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeature", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFeature indicates an expected call of CreateFeature.
func (mr *MockStoreMockRecorder) CreateFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeature", reflect.TypeOf((*MockStore)(nil).CreateFeature), ctx, opt)
}

// CreateGroup mocks base method.
func (m *MockStore) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockStoreMockRecorder) CreateGroup(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockStore)(nil).CreateGroup), ctx, opt)
}

// CreateRevision mocks base method.
func (m *MockStore) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRevision", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRevision indicates an expected call of CreateRevision.
func (mr *MockStoreMockRecorder) CreateRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRevision", reflect.TypeOf((*MockStore)(nil).CreateRevision), ctx, opt)
}

// GetCachedGroup mocks base method.
func (m *MockStore) GetCachedGroup(ctx context.Context, id int) (*types.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCachedGroup", ctx, id)
	ret0, _ := ret[0].(*types.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCachedGroup indicates an expected call of GetCachedGroup.
func (mr *MockStoreMockRecorder) GetCachedGroup(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCachedGroup", reflect.TypeOf((*MockStore)(nil).GetCachedGroup), ctx, id)
}

// GetEntity mocks base method.
func (m *MockStore) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", ctx, id)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntity indicates an expected call of GetEntity.
func (mr *MockStoreMockRecorder) GetEntity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockStore)(nil).GetEntity), ctx, id)
}

// GetEntityByName mocks base method.
func (m *MockStore) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByName", ctx, name)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByName indicates an expected call of GetEntityByName.
func (mr *MockStoreMockRecorder) GetEntityByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByName", reflect.TypeOf((*MockStore)(nil).GetEntityByName), ctx, name)
}

// GetFeature mocks base method.
func (m *MockStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeature", ctx, id)
	ret0, _ := ret[0].(*types.Feature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeature indicates an expected call of GetFeature.
func (mr *MockStoreMockRecorder) GetFeature(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeature", reflect.TypeOf((*MockStore)(nil).GetFeature), ctx, id)
}

// GetFeatureByName mocks base method.
func (m *MockStore) GetFeatureByName(ctx context.Context, groupName, featureName string) (*types.Feature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeatureByName", ctx, groupName, featureName)
	ret0, _ := ret[0].(*types.Feature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeatureByName indicates an expected call of GetFeatureByName.
func (mr *MockStoreMockRecorder) GetFeatureByName(ctx, groupName, featureName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeatureByName", reflect.TypeOf((*MockStore)(nil).GetFeatureByName), ctx, groupName, featureName)
}

// GetGroup mocks base method.
func (m *MockStore) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", ctx, id)
	ret0, _ := ret[0].(*types.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockStoreMockRecorder) GetGroup(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockStore)(nil).GetGroup), ctx, id)
}

// GetGroupByName mocks base method.
func (m *MockStore) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByName", ctx, name)
	ret0, _ := ret[0].(*types.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByName indicates an expected call of GetGroupByName.
func (mr *MockStoreMockRecorder) GetGroupByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByName", reflect.TypeOf((*MockStore)(nil).GetGroupByName), ctx, name)
}

// GetRevision mocks base method.
func (m *MockStore) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevision", ctx, id)
	ret0, _ := ret[0].(*types.Revision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevision indicates an expected call of GetRevision.
func (mr *MockStoreMockRecorder) GetRevision(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevision", reflect.TypeOf((*MockStore)(nil).GetRevision), ctx, id)
}

// GetRevisionBy mocks base method.
func (m *MockStore) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevisionBy", ctx, groupID, revision)
	ret0, _ := ret[0].(*types.Revision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevisionBy indicates an expected call of GetRevisionBy.
func (mr *MockStoreMockRecorder) GetRevisionBy(ctx, groupID, revision interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevisionBy", reflect.TypeOf((*MockStore)(nil).GetRevisionBy), ctx, groupID, revision)
}

// ListCachedFeature mocks base method.
func (m *MockStore) ListCachedFeature(ctx context.Context, opt metadata.ListCachedFeatureOpt) types.FeatureList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCachedFeature", ctx, opt)
	ret0, _ := ret[0].(types.FeatureList)
	return ret0
}

// ListCachedFeature indicates an expected call of ListCachedFeature.
func (mr *MockStoreMockRecorder) ListCachedFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCachedFeature", reflect.TypeOf((*MockStore)(nil).ListCachedFeature), ctx, opt)
}

// ListEntity mocks base method.
func (m *MockStore) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntity", ctx, entityIDs)
	ret0, _ := ret[0].(types.EntityList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntity indicates an expected call of ListEntity.
func (mr *MockStoreMockRecorder) ListEntity(ctx, entityIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntity", reflect.TypeOf((*MockStore)(nil).ListEntity), ctx, entityIDs)
}

// ListFeature mocks base method.
func (m *MockStore) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeature", ctx, opt)
	ret0, _ := ret[0].(types.FeatureList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeature indicates an expected call of ListFeature.
func (mr *MockStoreMockRecorder) ListFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeature", reflect.TypeOf((*MockStore)(nil).ListFeature), ctx, opt)
}

// ListGroup mocks base method.
func (m *MockStore) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGroup", ctx, entityID, groupIDs)
	ret0, _ := ret[0].(types.GroupList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGroup indicates an expected call of ListGroup.
func (mr *MockStoreMockRecorder) ListGroup(ctx, entityID, groupIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGroup", reflect.TypeOf((*MockStore)(nil).ListGroup), ctx, entityID, groupIDs)
}

// ListRevision mocks base method.
func (m *MockStore) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRevision", ctx, groupID)
	ret0, _ := ret[0].(types.RevisionList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRevision indicates an expected call of ListRevision.
func (mr *MockStoreMockRecorder) ListRevision(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRevision", reflect.TypeOf((*MockStore)(nil).ListRevision), ctx, groupID)
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

// Refresh mocks base method.
func (m *MockStore) Refresh() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockStoreMockRecorder) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockStore)(nil).Refresh))
}

// UpdateEntity mocks base method.
func (m *MockStore) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEntity", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEntity indicates an expected call of UpdateEntity.
func (mr *MockStoreMockRecorder) UpdateEntity(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEntity", reflect.TypeOf((*MockStore)(nil).UpdateEntity), ctx, opt)
}

// UpdateFeature mocks base method.
func (m *MockStore) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFeature", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFeature indicates an expected call of UpdateFeature.
func (mr *MockStoreMockRecorder) UpdateFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeature", reflect.TypeOf((*MockStore)(nil).UpdateFeature), ctx, opt)
}

// UpdateGroup mocks base method.
func (m *MockStore) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockStoreMockRecorder) UpdateGroup(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockStore)(nil).UpdateGroup), ctx, opt)
}

// UpdateRevision mocks base method.
func (m *MockStore) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRevision", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRevision indicates an expected call of UpdateRevision.
func (mr *MockStoreMockRecorder) UpdateRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRevision", reflect.TypeOf((*MockStore)(nil).UpdateRevision), ctx, opt)
}

// WithTransaction mocks base method.
func (m *MockStore) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MockStoreMockRecorder) WithTransaction(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*MockStore)(nil).WithTransaction), ctx, fn)
}

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

// CreateEntity mocks base method.
func (m *MockDBStore) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntity", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateEntity indicates an expected call of CreateEntity.
func (mr *MockDBStoreMockRecorder) CreateEntity(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntity", reflect.TypeOf((*MockDBStore)(nil).CreateEntity), ctx, opt)
}

// CreateFeature mocks base method.
func (m *MockDBStore) CreateFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeature", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateFeature indicates an expected call of CreateFeature.
func (mr *MockDBStoreMockRecorder) CreateFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeature", reflect.TypeOf((*MockDBStore)(nil).CreateFeature), ctx, opt)
}

// CreateGroup mocks base method.
func (m *MockDBStore) CreateGroup(ctx context.Context, opt metadata.CreateGroupOpt) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateGroup", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateGroup indicates an expected call of CreateGroup.
func (mr *MockDBStoreMockRecorder) CreateGroup(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateGroup", reflect.TypeOf((*MockDBStore)(nil).CreateGroup), ctx, opt)
}

// CreateRevision mocks base method.
func (m *MockDBStore) CreateRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRevision", ctx, opt)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(string)
	ret2, _ := ret[2].(error)
	return ret0, ret1, ret2
}

// CreateRevision indicates an expected call of CreateRevision.
func (mr *MockDBStoreMockRecorder) CreateRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRevision", reflect.TypeOf((*MockDBStore)(nil).CreateRevision), ctx, opt)
}

// GetEntity mocks base method.
func (m *MockDBStore) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", ctx, id)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntity indicates an expected call of GetEntity.
func (mr *MockDBStoreMockRecorder) GetEntity(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockDBStore)(nil).GetEntity), ctx, id)
}

// GetEntityByName mocks base method.
func (m *MockDBStore) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntityByName", ctx, name)
	ret0, _ := ret[0].(*types.Entity)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetEntityByName indicates an expected call of GetEntityByName.
func (mr *MockDBStoreMockRecorder) GetEntityByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntityByName", reflect.TypeOf((*MockDBStore)(nil).GetEntityByName), ctx, name)
}

// GetFeature mocks base method.
func (m *MockDBStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeature", ctx, id)
	ret0, _ := ret[0].(*types.Feature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeature indicates an expected call of GetFeature.
func (mr *MockDBStoreMockRecorder) GetFeature(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeature", reflect.TypeOf((*MockDBStore)(nil).GetFeature), ctx, id)
}

// GetFeatureByName mocks base method.
func (m *MockDBStore) GetFeatureByName(ctx context.Context, groupName, featureName string) (*types.Feature, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeatureByName", ctx, groupName, featureName)
	ret0, _ := ret[0].(*types.Feature)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeatureByName indicates an expected call of GetFeatureByName.
func (mr *MockDBStoreMockRecorder) GetFeatureByName(ctx, groupName, featureName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeatureByName", reflect.TypeOf((*MockDBStore)(nil).GetFeatureByName), ctx, groupName, featureName)
}

// GetGroup mocks base method.
func (m *MockDBStore) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroup", ctx, id)
	ret0, _ := ret[0].(*types.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroup indicates an expected call of GetGroup.
func (mr *MockDBStoreMockRecorder) GetGroup(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroup", reflect.TypeOf((*MockDBStore)(nil).GetGroup), ctx, id)
}

// GetGroupByName mocks base method.
func (m *MockDBStore) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetGroupByName", ctx, name)
	ret0, _ := ret[0].(*types.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetGroupByName indicates an expected call of GetGroupByName.
func (mr *MockDBStoreMockRecorder) GetGroupByName(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetGroupByName", reflect.TypeOf((*MockDBStore)(nil).GetGroupByName), ctx, name)
}

// GetRevision mocks base method.
func (m *MockDBStore) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevision", ctx, id)
	ret0, _ := ret[0].(*types.Revision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevision indicates an expected call of GetRevision.
func (mr *MockDBStoreMockRecorder) GetRevision(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevision", reflect.TypeOf((*MockDBStore)(nil).GetRevision), ctx, id)
}

// GetRevisionBy mocks base method.
func (m *MockDBStore) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevisionBy", ctx, groupID, revision)
	ret0, _ := ret[0].(*types.Revision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevisionBy indicates an expected call of GetRevisionBy.
func (mr *MockDBStoreMockRecorder) GetRevisionBy(ctx, groupID, revision interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevisionBy", reflect.TypeOf((*MockDBStore)(nil).GetRevisionBy), ctx, groupID, revision)
}

// ListEntity mocks base method.
func (m *MockDBStore) ListEntity(ctx context.Context, entityIDs *[]int) (types.EntityList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntity", ctx, entityIDs)
	ret0, _ := ret[0].(types.EntityList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListEntity indicates an expected call of ListEntity.
func (mr *MockDBStoreMockRecorder) ListEntity(ctx, entityIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntity", reflect.TypeOf((*MockDBStore)(nil).ListEntity), ctx, entityIDs)
}

// ListFeature mocks base method.
func (m *MockDBStore) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) (types.FeatureList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeature", ctx, opt)
	ret0, _ := ret[0].(types.FeatureList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListFeature indicates an expected call of ListFeature.
func (mr *MockDBStoreMockRecorder) ListFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeature", reflect.TypeOf((*MockDBStore)(nil).ListFeature), ctx, opt)
}

// ListGroup mocks base method.
func (m *MockDBStore) ListGroup(ctx context.Context, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListGroup", ctx, entityID, groupIDs)
	ret0, _ := ret[0].(types.GroupList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListGroup indicates an expected call of ListGroup.
func (mr *MockDBStoreMockRecorder) ListGroup(ctx, entityID, groupIDs interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListGroup", reflect.TypeOf((*MockDBStore)(nil).ListGroup), ctx, entityID, groupIDs)
}

// ListRevision mocks base method.
func (m *MockDBStore) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRevision", ctx, groupID)
	ret0, _ := ret[0].(types.RevisionList)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListRevision indicates an expected call of ListRevision.
func (mr *MockDBStoreMockRecorder) ListRevision(ctx, groupID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRevision", reflect.TypeOf((*MockDBStore)(nil).ListRevision), ctx, groupID)
}

// UpdateEntity mocks base method.
func (m *MockDBStore) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEntity", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateEntity indicates an expected call of UpdateEntity.
func (mr *MockDBStoreMockRecorder) UpdateEntity(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEntity", reflect.TypeOf((*MockDBStore)(nil).UpdateEntity), ctx, opt)
}

// UpdateFeature mocks base method.
func (m *MockDBStore) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFeature", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateFeature indicates an expected call of UpdateFeature.
func (mr *MockDBStoreMockRecorder) UpdateFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeature", reflect.TypeOf((*MockDBStore)(nil).UpdateFeature), ctx, opt)
}

// UpdateGroup mocks base method.
func (m *MockDBStore) UpdateGroup(ctx context.Context, opt metadata.UpdateGroupOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateGroup", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateGroup indicates an expected call of UpdateGroup.
func (mr *MockDBStoreMockRecorder) UpdateGroup(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateGroup", reflect.TypeOf((*MockDBStore)(nil).UpdateGroup), ctx, opt)
}

// UpdateRevision mocks base method.
func (m *MockDBStore) UpdateRevision(ctx context.Context, opt metadata.UpdateRevisionOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRevision", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// UpdateRevision indicates an expected call of UpdateRevision.
func (mr *MockDBStoreMockRecorder) UpdateRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRevision", reflect.TypeOf((*MockDBStore)(nil).UpdateRevision), ctx, opt)
}

// WithTransaction mocks base method.
func (m *MockDBStore) WithTransaction(ctx context.Context, fn func(context.Context, metadata.DBStore) error) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "WithTransaction", ctx, fn)
	ret0, _ := ret[0].(error)
	return ret0
}

// WithTransaction indicates an expected call of WithTransaction.
func (mr *MockDBStoreMockRecorder) WithTransaction(ctx, fn interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "WithTransaction", reflect.TypeOf((*MockDBStore)(nil).WithTransaction), ctx, fn)
}

// MockCacheStore is a mock of CacheStore interface.
type MockCacheStore struct {
	ctrl     *gomock.Controller
	recorder *MockCacheStoreMockRecorder
}

// MockCacheStoreMockRecorder is the mock recorder for MockCacheStore.
type MockCacheStoreMockRecorder struct {
	mock *MockCacheStore
}

// NewMockCacheStore creates a new mock instance.
func NewMockCacheStore(ctrl *gomock.Controller) *MockCacheStore {
	mock := &MockCacheStore{ctrl: ctrl}
	mock.recorder = &MockCacheStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockCacheStore) EXPECT() *MockCacheStoreMockRecorder {
	return m.recorder
}

// GetCachedGroup mocks base method.
func (m *MockCacheStore) GetCachedGroup(ctx context.Context, id int) (*types.Group, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCachedGroup", ctx, id)
	ret0, _ := ret[0].(*types.Group)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCachedGroup indicates an expected call of GetCachedGroup.
func (mr *MockCacheStoreMockRecorder) GetCachedGroup(ctx, id interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCachedGroup", reflect.TypeOf((*MockCacheStore)(nil).GetCachedGroup), ctx, id)
}

// ListCachedFeature mocks base method.
func (m *MockCacheStore) ListCachedFeature(ctx context.Context, opt metadata.ListCachedFeatureOpt) types.FeatureList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListCachedFeature", ctx, opt)
	ret0, _ := ret[0].(types.FeatureList)
	return ret0
}

// ListCachedFeature indicates an expected call of ListCachedFeature.
func (mr *MockCacheStoreMockRecorder) ListCachedFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListCachedFeature", reflect.TypeOf((*MockCacheStore)(nil).ListCachedFeature), ctx, opt)
}

// Refresh mocks base method.
func (m *MockCacheStore) Refresh() error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Refresh")
	ret0, _ := ret[0].(error)
	return ret0
}

// Refresh indicates an expected call of Refresh.
func (mr *MockCacheStoreMockRecorder) Refresh() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Refresh", reflect.TypeOf((*MockCacheStore)(nil).Refresh))
}

// MockSqlxContext is a mock of SqlxContext interface.
type MockSqlxContext struct {
	ctrl     *gomock.Controller
	recorder *MockSqlxContextMockRecorder
}

// MockSqlxContextMockRecorder is the mock recorder for MockSqlxContext.
type MockSqlxContextMockRecorder struct {
	mock *MockSqlxContext
}

// NewMockSqlxContext creates a new mock instance.
func NewMockSqlxContext(ctrl *gomock.Controller) *MockSqlxContext {
	mock := &MockSqlxContext{ctrl: ctrl}
	mock.recorder = &MockSqlxContextMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockSqlxContext) EXPECT() *MockSqlxContextMockRecorder {
	return m.recorder
}

// DriverName mocks base method.
func (m *MockSqlxContext) DriverName() string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DriverName")
	ret0, _ := ret[0].(string)
	return ret0
}

// DriverName indicates an expected call of DriverName.
func (mr *MockSqlxContextMockRecorder) DriverName() *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DriverName", reflect.TypeOf((*MockSqlxContext)(nil).DriverName))
}

// ExecContext mocks base method.
func (m *MockSqlxContext) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "ExecContext", varargs...)
	ret0, _ := ret[0].(sql.Result)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ExecContext indicates an expected call of ExecContext.
func (mr *MockSqlxContextMockRecorder) ExecContext(ctx, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ExecContext", reflect.TypeOf((*MockSqlxContext)(nil).ExecContext), varargs...)
}

// GetContext mocks base method.
func (m *MockSqlxContext) GetContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// GetContext indicates an expected call of GetContext.
func (mr *MockSqlxContextMockRecorder) GetContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetContext", reflect.TypeOf((*MockSqlxContext)(nil).GetContext), varargs...)
}

// Rebind mocks base method.
func (m *MockSqlxContext) Rebind(arg0 string) string {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Rebind", arg0)
	ret0, _ := ret[0].(string)
	return ret0
}

// Rebind indicates an expected call of Rebind.
func (mr *MockSqlxContextMockRecorder) Rebind(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Rebind", reflect.TypeOf((*MockSqlxContext)(nil).Rebind), arg0)
}

// SelectContext mocks base method.
func (m *MockSqlxContext) SelectContext(ctx context.Context, dest interface{}, query string, args ...interface{}) error {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, dest, query}
	for _, a := range args {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "SelectContext", varargs...)
	ret0, _ := ret[0].(error)
	return ret0
}

// SelectContext indicates an expected call of SelectContext.
func (mr *MockSqlxContextMockRecorder) SelectContext(ctx, dest, query interface{}, args ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, dest, query}, args...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SelectContext", reflect.TypeOf((*MockSqlxContext)(nil).SelectContext), varargs...)
}
