package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test"
)

func TestGetFeatureGroup(t *testing.T) {
	test.TestGetFeatureGroup(t, prepareStore)
}

func TestListFeatureGroup(t *testing.T) {
	test.TestListFeatureGroup(t, prepareStore)
}

func TestCreateFeatureGroup(t *testing.T) {
	test.TestCreateFeatureGroup(t, prepareStore)
}

func TestUpdateFeatureGroup(t *testing.T) {
	test.TestUpdateFeatureGroup(t, prepareStore)
}
