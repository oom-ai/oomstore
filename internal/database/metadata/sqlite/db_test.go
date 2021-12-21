package sqlite_test

import (
	"context"
	"os"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"

	"github.com/oom-ai/oomstore/internal/database/metadata/postgres"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"

	"github.com/oom-ai/oomstore/internal/database/metadata/sqlite"
)

func prepareStore(t *testing.T) (context.Context, metadata.Store) {
	return prepareDB(t)
}

func prepareDB(t *testing.T) (context.Context, *sqlite.DB) {
	ctx := context.Background()
	file, err := os.CreateTemp(t.TempDir(), "")
	if err != nil {
		t.Fatal(err)
	}
	file.Close()

	err = sqlite.CreateDatabase(ctx, types.SQLiteOpt{
		DBFile: file.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}

	store, err := sqlite.Open(ctx, &types.SQLiteOpt{
		DBFile: file.Name(),
	})
	if err != nil {
		t.Fatal(err)
	}
	return context.Background(), store
}

func TestCreateDatabase(t *testing.T) {
	ctx, db := prepareDB(t)
	defer db.Close()

	var tables []string
	err := db.SelectContext(ctx, &tables,
		`SELECT name
			FROM sqlite_master
			WHERE type = 'table'
			ORDER BY name;`)
	require.NoError(t, err)

	var wantTables []string
	for table := range postgres.META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range postgres.META_VIEW_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	wantTables = append(wantTables, "sqlite_sequence")

	assert.ElementsMatch(t, wantTables, tables)
}

func TestCreateEntity(t *testing.T) {
	test_impl.TestCreateEntity(t, prepareStore)
}

func TestGetEntity(t *testing.T) {
	test_impl.TestGetEntity(t, prepareStore)
}

func TestUpdateEntity(t *testing.T) {
	test_impl.TestUpdateEntity(t, prepareStore)
}

func TestListEntity(t *testing.T) {
	test_impl.TestListEntity(t, prepareStore)
}
