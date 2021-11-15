package postgres_test

import (
	"context"
	"sort"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/postgres"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareStore(t *testing.T) (context.Context, metadata.Store) {
	return prepareDB(t)
}

func prepareDB(t *testing.T) (context.Context, *postgres.DB) {
	ctx := context.Background()
	opt := runtime_pg.PostgresDbOpt
	pg, err := postgres.OpenDB(
		context.Background(),
		opt.Host,
		opt.Port,
		opt.User,
		opt.Password,
		"test",
	)
	require.NoError(t, err)
	_, err = pg.ExecContext(ctx, "drop database if exists oomstore")
	require.NoError(t, err)
	pg.Close()

	err = postgres.CreateDatabase(ctx, runtime_pg.PostgresDbOpt)
	require.NoError(t, err)

	db, err := postgres.Open(context.Background(), &runtime_pg.PostgresDbOpt)
	require.NoError(t, err)

	return ctx, db
}

func TestCreateDatabase(t *testing.T) {
	ctx, store := prepareDB(t)
	defer store.Close()

	var tables []string
	err := store.SelectContext(ctx, &tables,
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'public'
			ORDER BY table_name;`)
	require.NoError(t, err)

	var wantTables []string
	for table := range postgres.META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range postgres.META_VIEW_SCHEMAS {
		wantTables = append(wantTables, table)
	}

	sort.Slice(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	sort.Slice(wantTables, func(i, j int) bool {
		return wantTables[i] < wantTables[j]
	})
	assert.Equal(t, wantTables, tables)
}
