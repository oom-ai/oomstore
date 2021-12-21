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

func TestCreateGroup(t *testing.T) {
	test_impl.TestCreateGroup(t, prepareStore)
}

func TestUpdateGroup(t *testing.T) {
	test_impl.TestUpdateGroup(t, prepareStore)
}

func TestGetGroup(t *testing.T) {
	test_impl.TestGetGroup(t, prepareStore)
}

func TestListGroup(t *testing.T) {
	test_impl.TestListGroup(t, prepareStore)
}

func TestCreateFeature(t *testing.T) {
	test_impl.TestCreateFeature(t, prepareStore)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	test_impl.TestCreateFeatureWithSameName(t, prepareStore)
}

func TestCreateFeatureWithSQLKeyword(t *testing.T) {
	test_impl.TestCreateFeatureWithSQLKeyword(t, prepareStore)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	// According to SQLite Type Affinity, invalid data type will be regarded as Numeric,
	// so we cannot use `Create Table` to check whether data type is valid or not.
	// This issue will be auto-resolved when we infer DB value type from go value type.
	t.Skip()
	test_impl.TestCreateFeatureWithInvalidDataType(t, prepareStore)
}

func TestUpdateFeature(t *testing.T) {
	test_impl.TestUpdateFeature(t, prepareStore)
}

func TestGetFeature(t *testing.T) {
	test_impl.TestGetFeature(t, prepareStore)
}

func TestGetFeatureByName(t *testing.T) {
	test_impl.TestGetFeatureByName(t, prepareStore)
}

func TestListFeature(t *testing.T) {
	test_impl.TestListFeature(t, prepareStore)
}

func TestCacheListFeature(t *testing.T) {
	test_impl.TestCacheListFeature(t, prepareStore)
}

func TestCreateRevision(t *testing.T) {
	test_impl.TestCreateRevision(t, prepareStore)
}

func TestUpdateRevision(t *testing.T) {
	test_impl.TestUpdateRevision(t, prepareStore)
}

func TestGetRevision(t *testing.T) {
	test_impl.TestGetRevision(t, prepareStore)
}

func TestGetRevisionBy(t *testing.T) {
	test_impl.TestGetRevisionBy(t, prepareStore)
}

func TestListRevision(t *testing.T) {
	test_impl.TestListRevision(t, prepareStore)
}
