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

var DATABASE string

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	DATABASE = t.TempDir() + "/test.db"
	file, err := os.Create(DATABASE)
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

func destroyStore() {
	if err := os.RemoveAll(DATABASE); err != nil {
		panic(err)
	}
}

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

func TestSnapshot(t *testing.T) {
	test_impl.TestSnapshot(t, prepareStore, destroyStore)
}

func TestTableSchema(t *testing.T) {
	test_impl.TestTableSchema(t, prepareStore, destroyStore, func(ctx context.Context) {
		db, err := dbutil.OpenSQLite(DATABASE)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		if _, err = db.ExecContext(ctx, "create table `offline_batch_1_1`(`user` varchar(16), `age` smallint, `unix_milli` int)"); err != nil {
			t.Fatal(err)
		}
		if _, err = db.ExecContext(ctx, "insert into `offline_batch_1_1` VALUES ('1', 1, 1), ('2', 2, 100)"); err != nil {
			t.Fatal(err)
		}
	})
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, destroyStore)
}
