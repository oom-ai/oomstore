// Code generated by MockGen. DO NOT EDIT.
// Source: internal/database/metadatav2/store.go

// Package mock_metadatav2 is a generated GoMock package.
package mock_metadatav2

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2"
	types "github.com/oom-ai/oomstore/pkg/oomstore/types"
	typesv2 "github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
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

// BuildRevisionRanges mocks base method.
func (m *MockStore) BuildRevisionRanges(ctx context.Context, groupName string) []*types.RevisionRange {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "BuildRevisionRanges", ctx, groupName)
	ret0, _ := ret[0].([]*types.RevisionRange)
	return ret0
}

// BuildRevisionRanges indicates an expected call of BuildRevisionRanges.
func (mr *MockStoreMockRecorder) BuildRevisionRanges(ctx, groupName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "BuildRevisionRanges", reflect.TypeOf((*MockStore)(nil).BuildRevisionRanges), ctx, groupName)
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
func (m *MockStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateEntity", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateEntity indicates an expected call of CreateEntity.
func (mr *MockStoreMockRecorder) CreateEntity(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateEntity", reflect.TypeOf((*MockStore)(nil).CreateEntity), ctx, opt)
}

// CreateFeature mocks base method.
func (m *MockStore) CreateFeature(ctx context.Context, opt metadatav2.CreateFeatureOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeature", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFeature indicates an expected call of CreateFeature.
func (mr *MockStoreMockRecorder) CreateFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeature", reflect.TypeOf((*MockStore)(nil).CreateFeature), ctx, opt)
}

// CreateFeatureGroup mocks base method.
func (m *MockStore) CreateFeatureGroup(ctx context.Context, opt metadatav2.CreateFeatureGroupOpt) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateFeatureGroup", ctx, opt)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateFeatureGroup indicates an expected call of CreateFeatureGroup.
func (mr *MockStoreMockRecorder) CreateFeatureGroup(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateFeatureGroup", reflect.TypeOf((*MockStore)(nil).CreateFeatureGroup), ctx, opt)
}

// CreateRevision mocks base method.
func (m *MockStore) CreateRevision(ctx context.Context, opt metadatav2.CreateRevisionOpt) (*typesv2.Revision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateRevision", ctx, opt)
	ret0, _ := ret[0].(*typesv2.Revision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateRevision indicates an expected call of CreateRevision.
func (mr *MockStoreMockRecorder) CreateRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateRevision", reflect.TypeOf((*MockStore)(nil).CreateRevision), ctx, opt)
}

// GetEntity mocks base method.
func (m *MockStore) GetEntity(ctx context.Context, name string) *typesv2.Entity {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetEntity", ctx, name)
	ret0, _ := ret[0].(*typesv2.Entity)
	return ret0
}

// GetEntity indicates an expected call of GetEntity.
func (mr *MockStoreMockRecorder) GetEntity(ctx, name interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetEntity", reflect.TypeOf((*MockStore)(nil).GetEntity), ctx, name)
}

// GetFeature mocks base method.
func (m *MockStore) GetFeature(ctx context.Context, featureName string) *typesv2.Feature {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeature", ctx, featureName)
	ret0, _ := ret[0].(*typesv2.Feature)
	return ret0
}

// GetFeature indicates an expected call of GetFeature.
func (mr *MockStoreMockRecorder) GetFeature(ctx, featureName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeature", reflect.TypeOf((*MockStore)(nil).GetFeature), ctx, featureName)
}

// GetFeatureGroup mocks base method.
func (m *MockStore) GetFeatureGroup(ctx context.Context, groupName string) *typesv2.FeatureGroup {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeatureGroup", ctx, groupName)
	ret0, _ := ret[0].(*typesv2.FeatureGroup)
	return ret0
}

// GetFeatureGroup indicates an expected call of GetFeatureGroup.
func (mr *MockStoreMockRecorder) GetFeatureGroup(ctx, groupName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeatureGroup", reflect.TypeOf((*MockStore)(nil).GetFeatureGroup), ctx, groupName)
}

// GetLatestRevision mocks base method.
func (m *MockStore) GetLatestRevision(ctx context.Context, groupName string) *typesv2.Revision {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLatestRevision", ctx, groupName)
	ret0, _ := ret[0].(*typesv2.Revision)
	return ret0
}

// GetLatestRevision indicates an expected call of GetLatestRevision.
func (mr *MockStoreMockRecorder) GetLatestRevision(ctx, groupName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLatestRevision", reflect.TypeOf((*MockStore)(nil).GetLatestRevision), ctx, groupName)
}

// GetRevision mocks base method.
func (m *MockStore) GetRevision(ctx context.Context, opt metadatav2.GetRevisionOpt) (*typesv2.Revision, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetRevision", ctx, opt)
	ret0, _ := ret[0].(*typesv2.Revision)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetRevision indicates an expected call of GetRevision.
func (mr *MockStoreMockRecorder) GetRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetRevision", reflect.TypeOf((*MockStore)(nil).GetRevision), ctx, opt)
}

// ListEntity mocks base method.
func (m *MockStore) ListEntity(ctx context.Context) typesv2.EntityList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListEntity", ctx)
	ret0, _ := ret[0].(typesv2.EntityList)
	return ret0
}

// ListEntity indicates an expected call of ListEntity.
func (mr *MockStoreMockRecorder) ListEntity(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListEntity", reflect.TypeOf((*MockStore)(nil).ListEntity), ctx)
}

// ListFeature mocks base method.
func (m *MockStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) typesv2.FeatureList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeature", ctx, opt)
	ret0, _ := ret[0].(typesv2.FeatureList)
	return ret0
}

// ListFeature indicates an expected call of ListFeature.
func (mr *MockStoreMockRecorder) ListFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeature", reflect.TypeOf((*MockStore)(nil).ListFeature), ctx, opt)
}

// ListFeatureGroup mocks base method.
func (m *MockStore) ListFeatureGroup(ctx context.Context, entityName *string) typesv2.FeatureGroupList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListFeatureGroup", ctx, entityName)
	ret0, _ := ret[0].(typesv2.FeatureGroupList)
	return ret0
}

// ListFeatureGroup indicates an expected call of ListFeatureGroup.
func (mr *MockStoreMockRecorder) ListFeatureGroup(ctx, entityName interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListFeatureGroup", reflect.TypeOf((*MockStore)(nil).ListFeatureGroup), ctx, entityName)
}

// ListRevision mocks base method.
func (m *MockStore) ListRevision(ctx context.Context, opt metadatav2.ListRevisionOpt) typesv2.RevisionList {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListRevision", ctx, opt)
	ret0, _ := ret[0].(typesv2.RevisionList)
	return ret0
}

// ListRevision indicates an expected call of ListRevision.
func (mr *MockStoreMockRecorder) ListRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListRevision", reflect.TypeOf((*MockStore)(nil).ListRevision), ctx, opt)
}

// UpdateEntity mocks base method.
func (m *MockStore) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateEntity", ctx, opt)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateEntity indicates an expected call of UpdateEntity.
func (mr *MockStoreMockRecorder) UpdateEntity(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateEntity", reflect.TypeOf((*MockStore)(nil).UpdateEntity), ctx, opt)
}

// UpdateFeature mocks base method.
func (m *MockStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFeature", ctx, opt)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFeature indicates an expected call of UpdateFeature.
func (mr *MockStoreMockRecorder) UpdateFeature(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeature", reflect.TypeOf((*MockStore)(nil).UpdateFeature), ctx, opt)
}

// UpdateFeatureGroup mocks base method.
func (m *MockStore) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateFeatureGroup", ctx, opt)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateFeatureGroup indicates an expected call of UpdateFeatureGroup.
func (mr *MockStoreMockRecorder) UpdateFeatureGroup(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateFeatureGroup", reflect.TypeOf((*MockStore)(nil).UpdateFeatureGroup), ctx, opt)
}

// UpdateRevision mocks base method.
func (m *MockStore) UpdateRevision(ctx context.Context, opt metadatav2.UpdateRevisionOpt) (int64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateRevision", ctx, opt)
	ret0, _ := ret[0].(int64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateRevision indicates an expected call of UpdateRevision.
func (mr *MockStoreMockRecorder) UpdateRevision(ctx, opt interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateRevision", reflect.TypeOf((*MockStore)(nil).UpdateRevision), ctx, opt)
}
