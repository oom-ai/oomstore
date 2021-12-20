package mysql_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/mysql"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_mysql"
)

func prepareStore() (context.Context, offline.Store) {
	ctx, db := runtime_mysql.PrepareDB()

	if _, err := db.ExecContext(ctx, fmt.Sprintf("CREATE DATABASE %s", runtime_mysql.MySQLDbOpt.Database)); err != nil {
		panic(err)
	}
	db.Close()

	store, err := mysql.Open(&runtime_mysql.MySQLDbOpt)
	if err != nil {
		panic(err)
	}

	return ctx, store
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore)
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore)
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore)
}
