package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test"
)

func TestCreateFeature(t *testing.T) {
	test.TestCreateFeature(t, prepareStore)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	test.TestCreateFeatureWithSameName(t, prepareStore)
}

func TestCreateFeatureWithSQLKeywrod(t *testing.T) {
	test.TestCreateFeatureWithSQLKeywrod(t, prepareStore)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	test.TestCreateFeatureWithInvalidDataType(t, prepareStore)
}

func TestGetFeature(t *testing.T) {
	test.TestGetFeature(t, prepareStore)
}

func TestListFeature(t *testing.T) {
	test.TestListFeature(t, prepareStore)
}

func TestUpdateFeature(t *testing.T) {
	test.TestUpdateFeature(t, prepareStore)
}
