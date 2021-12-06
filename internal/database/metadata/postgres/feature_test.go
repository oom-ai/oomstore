package postgres_test

import (
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata/test_impl"
)

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
