package test_impl

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/offline"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestTableSchema(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn, prepareSampleTable func(ctx context.Context)) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	prepareSampleTable(ctx)

	actual, err := store.TableSchema(ctx, offline.TableSchemaOpt{
		TableName:      "offline_batch_1_1",
		CheckTimeRange: true,
	})
	require.NoError(t, err)
	require.Equal(t, 3, len(actual.Fields))

	expected := types.DataTableSchema{
		Fields: []types.DataTableFieldSchema{
			{
				Name:      "user",
				ValueType: types.String,
			},
			{
				Name:      "age",
				ValueType: types.Int64,
			},
			{
				Name:      "unix_milli",
				ValueType: types.Int64,
			},
		},
		TimeRange: types.DataTableTimeRange{
			MinUnixMilli: int64Ptr(1),
			MaxUnixMilli: int64Ptr(100),
		},
	}
	assert.ElementsMatch(t, expected.Fields, actual.Fields)
	assert.Equal(t, expected.TimeRange, actual.TimeRange)
}

func int64Ptr(i int64) *int64 {
	return &i
}
