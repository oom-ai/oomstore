package mysql_test

import (
	"context"
	"fmt"
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
	runtime_mysql.Reset(DATABASE)
}

func prepareStore(t *testing.T) (context.Context, online.Store) {
	ctx, db := runtime_mysql.PrepareDB(t, DATABASE)
	opt := runtime_mysql.GetOpt(DATABASE)

	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", opt.Database)); err != nil {
		t.Fatal(err)
	}
	db.Close()

	store, err := mysql.Open(opt)
	if err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

func TestOpen(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestOpen(t, prepareStore)
}

func TestGetExisted(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestGetExisted(t, prepareStore)
}

func TestGetNotExistedEntityKey(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestGetNotExistedEntityKey(t, prepareStore)
}

func TestMultiGet(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestMultiGet(t, prepareStore)
}

func TestPurgeRemovesSpecifiedRevision(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestPurgeRemovesSpecifiedRevision(t, prepareStore)
}

func TestPurgeNotRemovesOtherRevisions(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestPurgeNotRemovesOtherRevisions(t, prepareStore)
}

func TestPing(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestPing(t, prepareStore)
}
