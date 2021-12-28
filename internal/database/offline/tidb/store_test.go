package tidb_test

import (
	"context"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mysql"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_tidb"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	runtime_tidb.CreateDatabase(DATABASE)

	store, err := mysql.Open(runtime_tidb.GetOpt(DATABASE))
	if err != nil {
		t.Fatal(err)
	}

	return context.Background(), store
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}
