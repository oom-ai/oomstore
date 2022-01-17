package tidb_test

import (
	"context"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/mysql"
	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_tidb"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, metadata.Store) {
	ctx := context.Background()
	opt := runtime_tidb.GetOpt(DATABASE)

	if err := mysql.CreateDatabase(ctx, opt); err != nil {
		t.Fatal(err)
	}
	store, err := mysql.Open(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}
	return ctx, store
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateDatabase(t *testing.T) {
	t.Cleanup(runtime_tidb.DestroyStore(DATABASE))

	ctx := context.Background()
	opt := runtime_tidb.GetOpt(DATABASE)
	if err := mysql.CreateDatabase(ctx, opt); err != nil {
		t.Fatal(err)
	}
	store, err := mysql.Open(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var tables []string
	err = store.SelectContext(ctx, &tables,
		`SELECT table_name
			FROM information_schema.tables
			WHERE table_schema = ?
			ORDER BY table_name;`,
		opt.Database)
	require.NoError(t, err)

	var wantTables []string
	for table := range mysql.META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range mysql.META_VIEW_SCHEMAS {
		wantTables = append(wantTables, table)
	}

	assert.ElementsMatch(t, wantTables, tables)
}

func TestCreateEntity(t *testing.T) {
	test_impl.TestCreateEntity(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestGetEntity(t *testing.T) {
	test_impl.TestGetEntity(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestUpdateEntity(t *testing.T) {
	test_impl.TestUpdateEntity(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestListEntity(t *testing.T) {
	test_impl.TestListEntity(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestGetGroup(t *testing.T) {
	test_impl.TestGetGroup(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestListGroup(t *testing.T) {
	test_impl.TestListGroup(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateGroup(t *testing.T) {
	test_impl.TestCreateGroup(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestUpdateGroup(t *testing.T) {
	test_impl.TestUpdateGroup(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateFeature(t *testing.T) {
	test_impl.TestCreateFeature(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateFeatureWithSameName(t *testing.T) {
	test_impl.TestCreateFeatureWithSameName(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateFeatureWithSQLKeyword(t *testing.T) {
	test_impl.TestCreateFeatureWithSQLKeyword(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	test_impl.TestCreateFeatureWithInvalidDataType(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestGetFeature(t *testing.T) {
	test_impl.TestGetFeature(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestGetFeatureByName(t *testing.T) {
	test_impl.TestGetFeatureByName(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestListFeature(t *testing.T) {
	test_impl.TestListFeature(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCacheListFeature(t *testing.T) {
	test_impl.TestListCachedFeature(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestUpdateFeature(t *testing.T) {
	test_impl.TestUpdateFeature(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestCreateRevision(t *testing.T) {
	test_impl.TestCreateRevision(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestUpdateRevision(t *testing.T) {
	test_impl.TestUpdateRevision(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestGetRevision(t *testing.T) {
	test_impl.TestGetRevision(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestGetRevisionBy(t *testing.T) {
	test_impl.TestGetRevisionBy(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}

func TestListRevision(t *testing.T) {
	test_impl.TestListRevision(t, prepareStore, runtime_tidb.DestroyStore(DATABASE))
}
