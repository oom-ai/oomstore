package postgres_test

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestExport(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	// make test entities
	ctx := context.Background()
	entity := &types.Entity{
		Name:   "device",
		Length: 10,
	}
	features := types.FeatureList{
		{
			Name:        "model",
			DBValueType: "VARCHAR(32)",
		},
		{
			Name:        "price",
			DBValueType: "INT",
		},
	}

	// prepare test data
	schema := dbutil.BuildFeatureDataTableSchema("device_info_1", entity, features)
	_, err := db.ExecContext(ctx, schema)
	require.NoError(t, err)

	err = insertTestDataToBasic(db, ctx, "device_info_1", "1234", "xiaomi", 100)
	require.NoError(t, err)
	err = insertTestDataToBasic(db, ctx, "device_info_1", "1235", "apple", 200)
	require.NoError(t, err)
	err = insertTestDataToBasic(db, ctx, "device_info_1", "1236", "huawei", 300)
	require.NoError(t, err)
	err = insertTestDataToBasic(db, ctx, "device_info_1", "1237", "oneplus", 240)
	require.NoError(t, err)

	testCases := []struct {
		description string
		opt         offline.ExportOpt
		expected    [][]interface{}
	}{
		{
			description: "no features",
			opt: offline.ExportOpt{
				DataTable:    "device_info_1",
				EntityName:   "device",
				FeatureNames: []string{},
			},
			expected: [][]interface{}{{"1234"}, {"1235"}, {"1236"}, {"1237"}},
		},
		{
			description: "valid features and valid entity rows",
			opt: offline.ExportOpt{
				DataTable:    "device_info_1",
				EntityName:   "device",
				FeatureNames: []string{"model", "price"},
			},
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := db.Export(context.Background(), tc.opt)
			assert.NoError(t, err)
			values := make([][]interface{}, 0)
			for ele := range actual {
				values = append(values, ele.Record)
				assert.NoError(t, ele.Error)
			}
			assert.Equal(t, tc.expected, values)
		})
	}
}
