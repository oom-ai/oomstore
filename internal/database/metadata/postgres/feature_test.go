package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
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

func TestCacheGetFeature(t *testing.T) {
	test_impl.TestCacheGetFeature(t, prepareStore)
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

func TestCatchListFeature(t *testing.T) {
	test_impl.TestCatheListFeature(t, prepareStore)
}

func TestUpdateFeature(t *testing.T) {
	test_impl.TestUpdateFeature(t, prepareStore)
}
