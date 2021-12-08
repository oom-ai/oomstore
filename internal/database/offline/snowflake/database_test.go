package snowflake_test

import (
	"context"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/snowflake"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	return prepareDB(t)
}

func prepareDB(t *testing.T) (context.Context, *snowflake.DB) {
	// Have to hardcode because we have no tool to mock snowflake
	// The info are subject to change because each trial only lasts 30 days
	opt := types.SnowflakeOpt{
		Account:  "fka25816",
		User:     "yiksanchan",
		Password: "snowflakeYYDS1",
		Database: "test",
	}
	db, err := snowflake.Open(&opt)
	require.NoError(t, err)

	ctx := context.Background()
	query := "SELECT 1"
	rows, err := db.QueryContext(ctx, query) // no cancel is allowed
	require.NoError(t, err)

	defer rows.Close()
	var v int
	for rows.Next() {
		err := rows.Scan(&v)
		assert.NoError(t, err)
		assert.Equal(t, 1, v)
	}
	assert.NoError(t, rows.Err())
	return ctx, db
}

// func TestPing(t *testing.T) {
// 	test_impl.TestPing(t, prepareStore)
// }
