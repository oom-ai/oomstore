package runtime_mysql_test

import (
	"sort"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata/mysql"
	"github.com/ethhte88/oomstore/internal/database/test/runtime_mysql"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateDatabase(t *testing.T) {
	ctx, store := runtime_mysql.PrepareDB()
	defer store.Close()

	var tables []string
	err := store.SelectContext(ctx, &tables,
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = 'test'
			ORDER BY table_name;`)
	require.NoError(t, err)

	var wantTables []string
	for table := range mysql.META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range mysql.META_VIEW_SCHEMAS {
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
