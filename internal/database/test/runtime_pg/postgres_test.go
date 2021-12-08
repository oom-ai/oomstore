package runtime_pg_test

import (
	"sort"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata/postgres"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_pg"
	"github.com/stretchr/testify/require"
	"gotest.tools/assert"
)

func TestCreateDatabase(t *testing.T) {
	ctx, store := runtime_pg.PrepareDB()
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
