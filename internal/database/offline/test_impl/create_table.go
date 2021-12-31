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
	entity := types.Entity{
		Name:   "device",
		Length: 16,
	}
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
			describtion: "with unix milli",
			opt: offline.CreateTableOpt{
				TableName:      tableName,
				Entity:         &entity,
				Features:       features,
				WithUnixMillis: true,
			},
		},
		{
			describtion: "without milli",
			opt: offline.CreateTableOpt{
				TableName:      tableName,
				Entity:         &entity,
				Features:       features,
				WithUnixMillis: false,
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
