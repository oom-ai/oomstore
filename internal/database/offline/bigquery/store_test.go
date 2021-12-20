package bigquery_test

import (
	"context"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
	"testing"
)

func prepareStore() (context.Context, offline.Store) {
	ctx, db := prepareDB()
	return ctx, db
}

func prepareDB() (context.Context, *bigquery.DB) {
	ctx := context.Background()
	opt := types.BigQueryOpt{
		ProjectID: "oom-feature-store",
		DatasetID: "test",
	}
	db, err := bigquery.Open(ctx, &opt)
	if err != nil {
		panic(err)
	}
	return ctx, db
}

func TestPing(t *testing.T) {
	// skip this unit test until we can put credentials to env
	t.Skip()
	test_impl.TestPing(t, prepareStore)
}

func TestImport(t *testing.T) {
	t.Skip()
	test_impl.TestImport(t, prepareStore)

	ctx, db := prepareDB()
	table := db.Dataset("test").Table("offline_1_1")
	err := table.Delete(ctx)
	require.NoError(t, err)
}

func TestExport(t *testing.T) {
	t.Skip()
	test_impl.TestExport(t, prepareStore)

	ctx, db := prepareDB()
	table := db.Dataset("test").Table("offline_1_1")
	err := table.Delete(ctx)
	require.NoError(t, err)
}
