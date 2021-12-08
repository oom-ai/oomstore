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
	return prepareDB()
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
	return context.Background(), db
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}
