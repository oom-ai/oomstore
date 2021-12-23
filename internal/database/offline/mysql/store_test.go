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
	ctx := context.Background()
	opt := runtime_mysql.GetOpt(DATABASE)
	db, err := dbutil.OpenMysqlDB(
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		"",
	)
	if err != nil {
		t.Fatal(err)
	}

	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", DATABASE)); err != nil {
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
	test_impl.TestPing(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}
