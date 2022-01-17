package mysql_test

import (
	"context"
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
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	runtime_mysql.CreateDatabase(DATABASE)

	store, err := mysql.Open(runtime_mysql.GetOpt(DATABASE))
	if err != nil {
		t.Fatal(err)
	}

	return context.Background(), store
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

func TestSnapshot(t *testing.T) {
	test_impl.TestSnapshot(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestPush(t *testing.T) {
	test_impl.TestPush(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}

func TestTableSchema(t *testing.T) {
	test_impl.TestTableSchema(t, prepareStore, runtime_mysql.DestroyStore(DATABASE), func(ctx context.Context) {
		opt := runtime_mysql.GetOpt(DATABASE)
		db, err := dbutil.OpenMysqlDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		if _, err := db.ExecContext(ctx, "create table `offline_batch_1_1`(`user` varchar(16), `age` smallint)"); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, runtime_mysql.DestroyStore(DATABASE))
}
