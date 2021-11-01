package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/postgres"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
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
	schema = dbutil.BuildFeatureDataTableSchema("device_info_15", entity, features)
	_, err = db.ExecContext(ctx, schema)
	require.NoError(t, err)

	err = insertTestData(db, ctx, "device_info_1", "1234", "xiaomi", 100)
	require.NoError(t, err)
	err = insertTestData(db, ctx, "device_info_1", "1235", "apple", 200)
	require.NoError(t, err)
	err = insertTestData(db, ctx, "device_info_15", "1234", "huawei", 300)
	require.NoError(t, err)
	err = insertTestData(db, ctx, "device_info_15", "1235", "oneplus", 240)
	require.NoError(t, err)

	testCases := []struct {
		description string
		opt         offline.JoinOpt
		expected    map[string]dbutil.RowMap
	}{
		{
			description: "no features",
			opt: offline.JoinOpt{
				Features: make(types.FeatureList, 0),
			},
			expected: make(map[string]dbutil.RowMap),
		},
		{
			description: "no entity rows",
			opt: offline.JoinOpt{
				Features:   features,
				EntityRows: make([]types.EntityRow, 0),
			},
			expected: make(map[string]dbutil.RowMap),
		},
		{
			description: "valid features and valid entity rows",
			opt: offline.JoinOpt{
				Entity: entity,
				EntityRows: []types.EntityRow{
					{
						EntityKey: "1234",
						UnixTime:  10,
					},
					{
						EntityKey: "1234",
						UnixTime:  20,
					},
					{
						EntityKey: "1235",
						UnixTime:  5,
					},
					{
						EntityKey: "1235",
						UnixTime:  15,
					},
				},
				RevisionRanges: []*types.RevisionRange{
					{
						MinRevision: 1,
						MaxRevision: 15,
						DataTable:   "device_info_1",
					},
					{
						MinRevision: 15,
						MaxRevision: 25,
						DataTable:   "device_info_15",
					},
				},
				Features: features,
			},
			expected: map[string]dbutil.RowMap{
				"1234,10": {
					"entity_key": "1234",
					"unix_time":  int64(10),
					"model":      "xiaomi",
					"price":      int64(100),
				},
				"1234,20": {
					"entity_key": "1234",
					"unix_time":  int64(20),
					"model":      "huawei",
					"price":      int64(300),
				},
				"1235,5": {
					"entity_key": "1235",
					"unix_time":  int64(5),
					"model":      "apple",
					"price":      int64(200),
				},
				"1235,15": {
					"entity_key": "1235",
					"unix_time":  int64(15),
					"model":      "oneplus",
					"price":      int64(240),
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := db.Join(context.Background(), tc.opt)
			assert.NoError(t, err)
			assert.Equal(t, tc.expected, actual)
		})
	}
}

func insertTestData(db *postgres.DB, ctx context.Context, tableName string, device, model string, price int32) error {
	query := fmt.Sprintf("INSERT INTO %s(device, model, price) VALUES($1, $2, $3)", tableName)
	_, err := db.ExecContext(ctx, query, device, model, price)
	return err
}
