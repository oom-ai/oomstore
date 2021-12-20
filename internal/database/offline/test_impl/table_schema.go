package test_impl

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func TestTableSchema(t *testing.T, prepareStore PrepareStoreFn, prepareSampleTable func(ctx context.Context)) {
	ctx, store := prepareStore()
	defer store.Close()

	prepareSampleTable(ctx)

	actual, err := store.TableSchema(ctx, "user")
	require.NoError(t, err)
	require.Equal(t, 2, len(actual.Fields))

	expected := types.DataTableSchema{
		Fields: []types.DataTableFieldSchema{
			{
				Name:      "user",
				ValueType: types.STRING,
			},
			{
				Name:      "age",
				ValueType: types.INT64,
			},
		},
	}
	require.ElementsMatch(t, expected.Fields, actual.Fields)
}
