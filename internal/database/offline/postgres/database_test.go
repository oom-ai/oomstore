package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/postgres"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	return prepareDB(t)
}

func prepareDB(t *testing.T) (context.Context, *postgres.DB) {
	ctx := context.Background()
	opt := runtime_pg.PostgresDbOpt
	store, err := postgres.Open(&types.PostgresOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Password: opt.Password,
		Database: "test",
	})
	require.NoError(t, err)
	defer store.Close()

	_, err = store.ExecContext(context.Background(), fmt.Sprintf("DROP DATABASE IF EXISTS %s", opt.Database))
	require.NoError(t, err)

	_, err = store.ExecContext(context.Background(), fmt.Sprintf("CREATE DATABASE %s", opt.Database))
	require.NoError(t, err)

	db, err := postgres.Open(&runtime_pg.PostgresDbOpt)
	require.NoError(t, err)
	return ctx, db
}
