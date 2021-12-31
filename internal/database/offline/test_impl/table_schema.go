package test_impl

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestTableSchema(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn, prepareSampleTable func(ctx context.Context)) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	prepareSampleTable(ctx)

	actual, err := store.TableSchema(ctx, "offline_batch_1_1")
	require.NoError(t, err)
	require.Equal(t, 2, len(actual.Fields))

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
		},
	}
	require.ElementsMatch(t, expected.Fields, actual.Fields)
}
