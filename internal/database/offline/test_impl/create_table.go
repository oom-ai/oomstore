package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
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
		describtion   string
		opt           offline.CreateTableOpt
		expectedError error
	}{
		{
			describtion: "cdc table (with unix_milli)",
			opt: offline.CreateTableOpt{
				TableName: tableName,
				Entity:    &entity,
				Features:  features,
				IsCDC:     true,
			},
		},
		{
			describtion: "snapshot table (without_milli)",
			opt: offline.CreateTableOpt{
				TableName: tableName,
				Entity:    &entity,
				Features:  features,
				IsCDC:     false,
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.describtion, func(t *testing.T) {
			err := store.CreateTable(ctx, testCase.opt)
			require.NoError(t, err)
		})
	}
}
