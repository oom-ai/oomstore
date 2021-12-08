package snowflake_test

import (
	"context"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/snowflake"
	"github.com/ethhte88/oomstore/internal/database/offline/test_impl"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func prepareStore() (context.Context, offline.Store) {
	ctx, db := prepareDB()
	if _, err := db.ExecContext(ctx, "CREATE DATABASE test"); err != nil {
		panic(err)
	}
	return ctx, db
}

func prepareDB() (context.Context, *snowflake.DB) {
	// Have to hardcode because we have no tool to mock snowflake
	// The info are subject to change because each trial only lasts 30 days
	opt := types.SnowflakeOpt{
		Account:  "fka25816",
		User:     "yiksanchan",
		Password: "snowflakeYYDS1",
	}
	db, err := snowflake.Open(&opt)
	if err != nil {
		panic(err)
	}

	ctx := context.Background()
	if _, err = db.ExecContext(ctx, "DROP DATABASE IF EXISTS test"); err != nil {
		panic(err)
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
