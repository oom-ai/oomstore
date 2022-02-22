package bigquery_test

import (
	"context"
	"fmt"
	"os"
	"strings"
	"testing"

	bq "cloud.google.com/go/bigquery"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/bigquery"
	"github.com/oom-ai/oomstore/internal/database/offline/test_impl"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

var DATASET_ID string

func init() {
	DATASET_ID = strings.ToLower(dbutil.RandString(20))
}

func prepareStore(t *testing.T) (context.Context, offline.Store) {
	ctx, db := prepareDB(t)
	return ctx, db
}

func prepareDB(t *testing.T) (context.Context, *bigquery.DB) {
	ctx := context.Background()
	opt := types.BigQueryOpt{
		ProjectID:   "oom-feature-store",
		DatasetID:   DATASET_ID,
		Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
	}
	db, err := bigquery.Open(ctx, &opt)
	if err != nil {
		t.Fatal(err)
	}

	err = db.Dataset(DATASET_ID).Create(ctx, &bq.DatasetMetadata{
		Location: "US",
	})
	if err != nil {
		t.Fatal(err)
	}
	return ctx, db
}

func destroyStore(datasetID string) func() {
	return func() {
		opt := types.BigQueryOpt{
			ProjectID:   "oom-feature-store",
			DatasetID:   datasetID,
			Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		}
		db, err := bigquery.Open(context.Background(), &opt)
		if err != nil {
			panic(err)
		}
		defer db.Close()

		if err := db.Dataset(DATASET_ID).DeleteWithContents(context.Background()); err != nil {
			panic(err)
		}
	}
}

func TestPing(t *testing.T) {
	test_impl.TestPing(t, prepareStore, destroyStore(DATASET_ID))
}

func TestImport(t *testing.T) {
	test_impl.TestImport(t, prepareStore, destroyStore(DATASET_ID))
}

func TestExport(t *testing.T) {
	test_impl.TestExport(t, prepareStore, destroyStore(DATASET_ID))
}

func TestJoin(t *testing.T) {
	test_impl.TestJoin(t, prepareStore, destroyStore(DATASET_ID))
}

func TestTableSchema(t *testing.T) {
	test_impl.TestTableSchema(t, prepareStore, destroyStore(DATASET_ID), func(ctx context.Context) {
		opt := types.BigQueryOpt{
			ProjectID:   "oom-feature-store",
			DatasetID:   DATASET_ID,
			Credentials: os.Getenv("GOOGLE_APPLICATION_CREDENTIALS"),
		}
		db, err := bigquery.Open(ctx, &opt)
		if err != nil {
			t.Fatal(err)
		}

		query := fmt.Sprintf("CREATE TABLE %s.offline_batch_1_1(`user` STRING, `age` BIGINT, `unix_milli` BIGINT)", DATASET_ID)
		if _, err = db.Query(query).Read(ctx); err != nil {
			t.Fatal(err)
		}
		if _, err = db.Query(fmt.Sprintf("insert into %s.offline_batch_1_1 VALUES ('1', 1, 1), ('2', 2, 100)", DATASET_ID)).Read(ctx); err != nil {
			t.Fatal(err)
		}
	})
}

func TestSnapshot(t *testing.T) {
	t.Skip()
	test_impl.TestSnapshot(t, prepareStore, destroyStore(DATASET_ID))
}

func TestCreateTable(t *testing.T) {
	test_impl.TestCreateTable(t, prepareStore, destroyStore(DATASET_ID))
}
