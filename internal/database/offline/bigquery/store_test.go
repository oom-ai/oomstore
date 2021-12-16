package bigquery_test

import (
	"context"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/internal/database/offline/bigquery"
	"github.com/ethhte88/oomstore/internal/database/offline/test_impl"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

func prepareStore() (context.Context, offline.Store) {
	ctx, db := prepareDB()
	return ctx, db
}

func prepareDB() (context.Context, *bigquery.DB) {
	ctx := context.Background()
	opt := types.BigQueryOpt{
		ProjectID: "oom-feature-store",
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
