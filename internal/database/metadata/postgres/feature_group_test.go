package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
)

func TestGetFeatureGroup(t *testing.T) {
	test_impl.TestGetFeatureGroup(t, prepareStore)
}

func TestListFeatureGroup(t *testing.T) {
	test_impl.TestListFeatureGroup(t, prepareStore)
}

func TestCreateFeatureGroup(t *testing.T) {
	test_impl.TestCreateFeatureGroup(t, prepareStore)
}

func TestUpdateFeatureGroup(t *testing.T) {
	test_impl.TestUpdateFeatureGroup(t, prepareStore)
}
