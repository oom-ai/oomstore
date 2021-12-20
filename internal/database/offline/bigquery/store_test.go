package bigquery_test

import (
	"context"
	"os"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := prepareDB(t)
	return ctx, db
}

func prepareDB(t *testing.T) (context.Context, *bigquery.DB) {
	ctx := context.Background()
	opt := types.BigQueryOpt{
		ProjectID:   "oom-feature-store",
		DatasetID:   "test",
		Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}
	db, err := bigquery.Open(ctx, &opt)
	if err != nil {
		t.Fatal(err)
	}
	return ctx, db
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore)
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore)

	ctx, db := prepareDB(t)
	table := db.Dataset("test").Table("offline_1_1")
	err := table.Delete(ctx)
	require.NoError(t, err)
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore)

	ctx, db := prepareDB(t)
	table := db.Dataset("test").Table("offline_1_1")
	err := table.Delete(ctx)
	require.NoError(t, err)
}
