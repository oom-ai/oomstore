package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/postgres"
	"github.com/ethhte88/oomstore/internal/database/offline/test_impl"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_pg"
)

func prepareStore() (context.Context, offline.Store) {
	ctx, db := runtime_pg.PrepareDB()

	_, err := db.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", runtime_pg.PostgresDbOpt.Database))
	if err != nil {
		panic(err)
	}
	db.Close()

	store, err := postgres.Open(&runtime_pg.PostgresDbOpt)
	if err != nil {
		panic(err)
	}

	return ctx, store
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore)
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore)
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore)
}

func TestTableSchema(t *testing.T) {
	test_impl.TestTableSchema(t, prepareStore, func(ctx context.Context) {
		opt := runtime_pg.PostgresDbOpt
		db, err := dbutil.OpenPostgresDB(opt.Host, opt.Port, opt.User, opt.Password, opt.Database)
		if err != nil {
			panic(err)
		}
		defer db.Close()
		if _, err := db.ExecContext(ctx, `create table "user"("user" varchar(16), "age" smallint)`); err != nil {
			panic(err)
		}
	})
}
