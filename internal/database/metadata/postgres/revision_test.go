package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata/test"
)

func TestCreateRevision(t *testing.T) {
	test.TestCreateRevision(t, prepareStore)
}

func TestUpdateRevision(t *testing.T) {
	test.TestUpdateRevision(t, prepareStore)
}

func TestGetRevision(t *testing.T) {
	test.TestGetRevision(t, prepareStore)
}

func TestGetRevisionBy(t *testing.T) {
	test.TestGetRevisionBy(t, prepareStore)
}

func TestListRevision(t *testing.T) {
	test.TestListRevision(t, prepareStore)
}
