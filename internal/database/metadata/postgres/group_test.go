package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
)

func TestCacheGetGroup(t *testing.T) {
	test_impl.TestCacheGetGroup(t, prepareStore)
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
