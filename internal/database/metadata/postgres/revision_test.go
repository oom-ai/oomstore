package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
)

func TestCreateRevision(t *testing.T) {
	test_impl.TestCreateRevision(t, prepareStore)
}

func TestUpdateRevision(t *testing.T) {
	test_impl.TestUpdateRevision(t, prepareStore)
}

func TestCacheGetRevision(t *testing.T) {
	test_impl.TestCacheGetRevision(t, prepareStore)
}

func TestCacheGetRevisionBy(t *testing.T) {
	test_impl.TestCacheGetRevisionBy(t, prepareStore)
}

func TestCacheListRevision(t *testing.T) {
	test_impl.TestCacheListRevision(t, prepareStore)
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
