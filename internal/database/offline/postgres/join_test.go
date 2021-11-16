package postgres_test

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
)

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore)
}
