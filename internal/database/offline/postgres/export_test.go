package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
)

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore)
}
