package mysql_test

import (
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata/test_impl"
)

func TestCreateEntity(t *testing.T) {
	test_impl.TestCreateEntity(t, prepareStore)
}

func TestGetEntity(t *testing.T) {
	test_impl.TestGetEntity(t, prepareStore)
}

func TestUpdateEntity(t *testing.T) {
	test_impl.TestUpdateEntity(t, prepareStore)
}

func TestListEntity(t *testing.T) {
	test_impl.TestListEntity(t, prepareStore)
}

func TestGetGroup(t *testing.T) {
	test_impl.TestGetGroup(t, prepareStore)
}

func TestListGroup(t *testing.T) {
	test_impl.TestListGroup(t, prepareStore)
}

func TestCreateGroup(t *testing.T) {
	test_impl.TestCreateGroup(t, prepareStore)
}

func TestUpdateGroup(t *testing.T) {
	test_impl.TestUpdateGroup(t, prepareStore)
}

func TestCreateFeature(t *testing.T) {
	test_impl.TestCreateFeature(t, prepareStore)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	test_impl.TestCreateFeatureWithSameName(t, prepareStore)
}

func TestCreateFeatureWithSQLKeyword(t *testing.T) {
	test_impl.TestCreateFeatureWithSQLKeyword(t, prepareStore)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	test_impl.TestCreateFeatureWithInvalidDataType(t, prepareStore)
}

func TestGetFeature(t *testing.T) {
	test_impl.TestGetFeature(t, prepareStore)
}

func TestGetFeatureByName(t *testing.T) {
	test_impl.TestGetFeatureByName(t, prepareStore)
}

func TestListFeature(t *testing.T) {
	test_impl.TestListFeature(t, prepareStore)
}

func TestCacheListFeature(t *testing.T) {
	test_impl.TestCacheListFeature(t, prepareStore)
}

func TestUpdateFeature(t *testing.T) {
	test_impl.TestUpdateFeature(t, prepareStore)
}

func TestCreateRevision(t *testing.T) {
	test_impl.TestCreateRevision(t, prepareStore)
}

func TestUpdateRevision(t *testing.T) {
	test_impl.TestUpdateRevision(t, prepareStore)
}

func TestGetRevision(t *testing.T) {
	test_impl.TestGetRevision(t, prepareStore)
}

func TestGetRevisionBy(t *testing.T) {
	test_impl.TestGetRevisionBy(t, prepareStore)
}

func TestListRevision(t *testing.T) {
	test_impl.TestListRevision(t, prepareStore)
}
