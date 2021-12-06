package postgres_test

import (
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata/test_impl"
)

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
