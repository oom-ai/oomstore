package mysql_test

import (
	"context"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/internal/database/online/mysql"
	"github.com/oom-ai/oomstore/internal/database/online/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_mysql"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, online.Store) {
	runtime_mysql.CreateDatabase(DATABASE)

	store, err := mysql.Open(runtime_mysql.GetOpt(DATABASE))
	if err != nil {
		t.Fatal(err)
	}

	return context.Background(), store
}

func TestOpen(t *testing.T) {
	test_impl.TestOpen(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestGetExisted(t *testing.T) {
	test_impl.TestGetExisted(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestGetNoRevision(t *testing.T) {
	test_impl.TestGetNoRevision(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestGetNotExistedEntityKey(t *testing.T) {
	test_impl.TestGetNotExistedEntityKey(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestMultiGet(t *testing.T) {
	test_impl.TestMultiGet(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestPush(t *testing.T) {
	test_impl.TestPush(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}
