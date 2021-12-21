package snowflake_test

import (
	"context"
	"os"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/snowflake"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := prepareDB(t)
	if _, err := db.ExecContext(ctx, "CREATE DATABASE test"); err != nil {
		t.Fatal(t)
	}
	return ctx, db
}

func prepareDB(t *testing.T) (context.Context, *snowflake.DB) {
	opt := types.SnowflakeOpt{
		Account:  os.Getenv("SNOWFLAKE_TEST_ACCOUNT"),
		User:     os.Getenv("SNOWFLAKE_TEST_USER"),
		Password: os.Getenv("SNOWFLAKE_TEST_PASSWORD"),
	}
	db, err := snowflake.Open(&opt)
	if err != nil {
		t.Fatal(err)
	}

	ctx := context.Background()
	if _, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS test"); err != nil {
		t.Fatal(err)
	}
	return ctx, db
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
