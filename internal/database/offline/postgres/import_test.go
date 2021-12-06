package postgres_test

import (
	"testing"

	"github.com/ethhte88/oomstore/internal/database/offline/test_impl"
)

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore)
}
