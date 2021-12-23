package sqlite_test

import (
	"context"
	"os"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/sqlite"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
)

func prepareStore(t *testing.T) (context.Context, online.Store) {
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

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore, destroyStore)
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore, destroyStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore, destroyStore)
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore, destroyStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, destroyStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, destroyStore)
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore)
}
