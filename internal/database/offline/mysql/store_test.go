package mysql_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mysql"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_mysql"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
	runtime_mysql.Reset(DATABASE)
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
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

func TestPing(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestPing(t, prepareStore)
}

func TestExport(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestExport(t, prepareStore)
}

func TestImport(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestImport(t, prepareStore)
}

func TestJoin(t *testing.T) {
	t.Cleanup(func() { runtime_mysql.Reset(DATABASE) })
	test_impl.TestJoin(t, prepareStore)
}
