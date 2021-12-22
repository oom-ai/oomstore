package postgres_test

import (
	"context"
	"fmt"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/postgres"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
	if err := runtime_pg.Reset(DATABASE); err != nil {
		panic(err)
	}
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := runtime_pg.PrepareDB(t, DATABASE)
	opt := runtime_pg.GetOpt(DATABASE)

	_, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", opt.Database))
	if err != nil {
		t.Fatal(err)
	}
	db.Close()

	store, err := postgres.Open(opt)
	if err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

func TestPing(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestPing(t, prepareStore)
}

func TestImport(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestImport(t, prepareStore)
}

func TestExport(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestExport(t, prepareStore)
}

func TestJoin(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestJoin(t, prepareStore)
}

func TestTableSchema(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestTableSchema(t, prepareStore, func(ctx context.Context) {
		opt := runtime_pg.GetOpt(DATABASE)
		db, err := dbutil.OpenPostgresDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
		if err != nil {
			t.Fatal(err)
		}
		defer db.Close()
		if _, err := db.ExecContext(ctx, `create table "user"("user" varchar(16), "age" smallint)`); err != nil {
			t.Fatal(err)
		}
	})
}
