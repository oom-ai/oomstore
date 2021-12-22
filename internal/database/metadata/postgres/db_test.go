package postgres_test

import (
	"context"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/postgres"
	"github.com/oom-ai/oomstore/internal/database/metadata/test_impl"
	"github.com/oom-ai/oomstore/internal/database/test/runtime_pg"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

var DATABASE string

func init() {
	DATABASE = strings.ToLower(dbutil.RandString(20))
	if err := runtime_pg.Reset(DATABASE); err != nil {
		panic(err)
	}
}

func prepareStore(t *testing.T) (context.Context, metadata.Store) {
	ctx, db := runtime_pg.PrepareDB(t, DATABASE)
	opt := runtime_pg.GetOpt(DATABASE)
	db.Close()

	if err := postgres.CreateDatabase(ctx, *opt); err != nil {
		t.Fatal(err)
	}
	store, err := postgres.Open(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}

	return ctx, store
}

func TestCreateDatabase(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })

	ctx, db := runtime_pg.PrepareDB(t, DATABASE)
	opt := runtime_pg.GetOpt(DATABASE)
	db.Close()

	if err := postgres.CreateDatabase(ctx, *opt); err != nil {
		t.Fatal(err)
	}
	store, err := postgres.Open(ctx, opt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var tables []string
	err = store.SelectContext(ctx, &tables,
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

	assert.ElementsMatch(t, wantTables, tables)
}

func TestCreateEntity(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateEntity(t, prepareStore)
}

func TestGetEntity(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetEntity(t, prepareStore)
}

func TestUpdateEntity(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestUpdateEntity(t, prepareStore)
}

func TestListEntity(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestListEntity(t, prepareStore)
}

func TestCreateFeature(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateFeature(t, prepareStore)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateFeatureWithSameName(t, prepareStore)
}

func TestCreateFeatureWithSQLKeyword(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateFeatureWithSQLKeyword(t, prepareStore)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateFeatureWithInvalidDataType(t, prepareStore)
}

func TestGetFeature(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetFeature(t, prepareStore)
}

func TestGetFeatureByName(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetFeatureByName(t, prepareStore)
}

func TestListFeature(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestListFeature(t, prepareStore)
}

func TestCacheListFeature(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCacheListFeature(t, prepareStore)
}

func TestUpdateFeature(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestUpdateFeature(t, prepareStore)
}

func TestGetGroup(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetGroup(t, prepareStore)
}

func TestListGroup(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestListGroup(t, prepareStore)
}

func TestCreateGroup(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateGroup(t, prepareStore)
}

func TestUpdateGroup(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestUpdateGroup(t, prepareStore)
}

func TestCreateRevision(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestCreateRevision(t, prepareStore)
}

func TestUpdateRevision(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestUpdateRevision(t, prepareStore)
}

func TestGetRevision(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetRevision(t, prepareStore)
}

func TestGetRevisionBy(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestGetRevisionBy(t, prepareStore)
}

func TestListRevision(t *testing.T) {
	t.Cleanup(func() { _ = runtime_pg.Reset(DATABASE) })
	test_impl.TestListRevision(t, prepareStore)
}
