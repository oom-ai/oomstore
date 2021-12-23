package sqlite_test

import (
	"context"
	"os"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/sqlite"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	file, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	conn, err := dbutil.OpenSQLite(file.Name())
	if err != nil {
		t.Fatal(err)
	}

	return context.Background(), &sqlite.DB{conn}
}

func destroyStore() {}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore)
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore, destroyStore)
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore, destroyStore)
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore, destroyStore)
}
