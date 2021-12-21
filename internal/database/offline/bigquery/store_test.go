package bigquery_test

import (
	"context"
	"os"
	"testing"

	bq "cloud.google.com/go/bigquery"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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

	_ = db.Dataset("test").DeleteWithContents(ctx)
	err = db.Dataset("test").Create(ctx, &bq.DatasetMetadata{
		Location: "US",
	})
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
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore)
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore)
}

func TestTableSchema(t *testing.T) {
	test_impl.TestTableSchema(t, prepareStore, func(ctx context.Context) {
		opt := types.BigQueryOpt{
			ProjectID:   "oom-feature-store",
			DatasetID:   "test",
			Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		}
		db, err := bigquery.Open(ctx, &opt)
		if err != nil {
			t.Fatal(err)
		}
		if _, err := db.Query("CREATE TABLE test.user(`user` STRING, `age` BIGINT)").Read(ctx); err != nil {
			t.Fatal(err)
		}
	})
}
