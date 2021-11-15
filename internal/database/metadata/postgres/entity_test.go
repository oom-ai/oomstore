package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test"
)

func TestCreateEntity(t *testing.T) {
	test.TestCreateEntity(t, prepareStore)
}

func TestGetEntity(t *testing.T) {
	test.TestGetEntity(t, prepareStore)
}

func TestUpdateEntity(t *testing.T) {
	test.TestUpdateEntity(t, prepareStore)
}

func TestListEntity(t *testing.T) {
	test.TestListEntity(t, prepareStore)
}
