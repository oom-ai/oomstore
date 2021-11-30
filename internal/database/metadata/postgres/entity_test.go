package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
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
