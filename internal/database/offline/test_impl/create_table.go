package test_impl

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestCreateTable(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	tableName := "new_table"
	entity := types.Entity{Name: "device"}
	features := []*types.Feature{
		{
			Name:      "price",
			ValueType: types.Int64,
		},
		{
			Name:      "model",
			ValueType: types.String,
		},
	}

	testCases := []struct {
		description   string
		opt           offline.CreateTableOpt
		expectedError error
	}{
		{
			description: "cdc table (with unix_milli)",
			opt: offline.CreateTableOpt{
				TableName:  tableName + "_0",
				EntityName: entity.Name,
				Features:   features,
				TableType:  types.TableStreamCdc,
			},
		},
		{
			description: "stream snapshot table (with unix_milli)",
			opt: offline.CreateTableOpt{
				TableName:  tableName + "_1",
				EntityName: entity.Name,
				Features:   features,
				TableType:  types.TableStreamSnapshot,
			},
		},
		{
			description: "batch snapshot table (without unix_milli)",
			opt: offline.CreateTableOpt{
				TableName:  tableName + "_2",
				EntityName: entity.Name,
				Features:   features,
				TableType:  types.TableBatchSnapshot,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.description, func(t *testing.T) {
			err := store.CreateTable(ctx, testCase.opt)
			require.NoError(t, err)
		})
	}
}
