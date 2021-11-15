package postgres_test

import (
	"context"
	"fmt"
	"math"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/offline/postgres"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestJoin(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()
	ctx := context.Background()

	// prepare test data
	entity := &typesv2.Entity{
		Name:   "device",
		Length: 10,
	}
	oneGroupFeatures, oneGroupFeatureMap := prepareFeatures(true)
	_, twoGroupFeatureMap := prepareFeatures(false)
	prepareTestData(ctx, db, t, entity, oneGroupFeatures, twoGroupFeatureMap)

	testCases := []struct {
		description string
		opt         offline.JoinOpt
		expected    *types.JoinResult
	}{
		{
			description: "no features",
			opt: offline.JoinOpt{
				FeatureMap: make(map[string]typesv2.FeatureList),
			},
			expected: nil,
		},
		{
			description: "no entity rows",
			opt: offline.JoinOpt{
				Entity:     *entity,
				FeatureMap: oneGroupFeatureMap,
				EntityRows: prepareEntityRows(true),
			},
			expected: nil,
		},
		{
			description: "one feature group",
			opt: offline.JoinOpt{
				Entity:           *entity,
				EntityRows:       prepareEntityRows(false),
				FeatureMap:       oneGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(true),
			},
			expected: prepareResult(true),
		},
		{
			description: "two feature groups",
			opt: offline.JoinOpt{
				Entity:           *entity,
				EntityRows:       prepareEntityRows(false),
				FeatureMap:       twoGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(false),
			},
			expected: prepareResult(false),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := db.Join(context.Background(), tc.opt)
			assert.NoError(t, err)
			if tc.expected == nil {
				assert.Equal(t, tc.expected, actual)
			} else {
				expectedValues := extractValues(tc.expected.Data)
				actualValues := extractValues(actual.Data)
				assert.ObjectsAreEqual(tc.expected.Header, actual.Header)
				assert.ObjectsAreEqual(expectedValues, actualValues)
			}
		})
	}
}

func prepareTestData(ctx context.Context, db *postgres.DB, t *testing.T, entity *typesv2.Entity, oneGroupFeatures typesv2.FeatureList, twoGroupFeatureMap map[string]typesv2.FeatureList) {
	schema := dbutil.BuildFeatureDataTableSchema("device_basic_1", entity, oneGroupFeatures)
	_, err := db.ExecContext(ctx, schema)
	require.NoError(t, err)
	schema = dbutil.BuildFeatureDataTableSchema("device_basic_15", entity, oneGroupFeatures)
	_, err = db.ExecContext(ctx, schema)
	require.NoError(t, err)
	schema = dbutil.BuildFeatureDataTableSchema("device_advanced_5", entity, twoGroupFeatureMap["device_advanced"])
	_, err = db.ExecContext(ctx, schema)
	require.NoError(t, err)

	err = insertTestDataToBasic(db, ctx, "device_basic_1", "1234", "xiaomi", 100)
	require.NoError(t, err)
	err = insertTestDataToBasic(db, ctx, "device_basic_1", "1235", "apple", 200)
	require.NoError(t, err)
	err = insertTestDataToBasic(db, ctx, "device_basic_15", "1234", "galaxy", 300)
	require.NoError(t, err)
	err = insertTestDataToBasic(db, ctx, "device_basic_15", "1235", "oneplus", 240)
	require.NoError(t, err)
	err = insertTestDataToAdvanced(db, ctx, "device_advanced_5", "1234", true)
	require.NoError(t, err)
	err = insertTestDataToAdvanced(db, ctx, "device_advanced_5", "1235", false)
	require.NoError(t, err)
}

func insertTestDataToBasic(db *postgres.DB, ctx context.Context, tableName string, device, model string, price int32) error {
	query := fmt.Sprintf("INSERT INTO %s(device, model, price) VALUES($1, $2, $3)", tableName)
	_, err := db.ExecContext(ctx, query, device, model, price)
	return err
}

func insertTestDataToAdvanced(db *postgres.DB, ctx context.Context, tableName string, device string, isActive bool) error {
	query := fmt.Sprintf("INSERT INTO %s(device, is_active) VALUES($1, $2)", tableName)
	_, err := db.ExecContext(ctx, query, device, isActive)
	return err
}

func prepareFeatures(oneGroup bool) (typesv2.FeatureList, map[string]typesv2.FeatureList) {
	price := &typesv2.Feature{
		Name:        "price",
		DBValueType: "INT",
		GroupID:     1,
	}
	model := &typesv2.Feature{
		Name:        "model",
		DBValueType: "VARCHAR(32)",
		GroupID:     1,
	}
	isActive := &typesv2.Feature{
		Name:        "is_active",
		DBValueType: "boolean",
		GroupID:     2,
	}

	if oneGroup {
		features := typesv2.FeatureList{model, price}
		featureMap := map[string]typesv2.FeatureList{
			"device_basic": {
				model, price,
			},
		}
		return features, featureMap
	} else {
		features := typesv2.FeatureList{model, price, isActive}
		featureMap := map[string]typesv2.FeatureList{
			"device_basic": {
				model, price,
			},
			"device_advanced": {isActive},
		}
		return features, featureMap
	}
}

func prepareRevisionRanges(oneGroup bool) map[string][]*metadata.RevisionRange {
	basic := []*metadata.RevisionRange{
		{
			MinRevision: 1,
			MaxRevision: 15,
			DataTable:   "device_basic_1",
		},
		{
			MinRevision: 15,
			MaxRevision: 25,
			DataTable:   "device_basic_15",
		},
	}
	advanced := []*metadata.RevisionRange{
		{
			MinRevision: 5,
			MaxRevision: math.MaxInt64,
			DataTable:   "device_advanced_5",
		},
	}
	if oneGroup {
		return map[string][]*metadata.RevisionRange{
			"device_basic": basic,
		}
	}
	return map[string][]*metadata.RevisionRange{
		"device_basic":    basic,
		"device_advanced": advanced,
	}
}

func prepareEntityRows(isEmpty bool) <-chan types.EntityRow {
	entityRows := make(chan types.EntityRow)
	go func() {
		defer close(entityRows)
		if isEmpty {
			return
		}
		entityRows <- types.EntityRow{
			EntityKey: "1234",
			UnixTime:  10,
		}
		entityRows <- types.EntityRow{
			EntityKey: "1234",
			UnixTime:  20,
		}
		entityRows <- types.EntityRow{
			EntityKey: "1235",
			UnixTime:  5,
		}
		entityRows <- types.EntityRow{
			EntityKey: "1235",
			UnixTime:  14,
		}
	}()
	return entityRows
}

func prepareResult(oneGroup bool) *types.JoinResult {
	header := []string{"entity_key", "unix_time", "model", "price", "is_active"}
	values := []map[string]interface{}{
		{
			"entity_key": "1234",
			"unix_time":  int64(10),
			"model":      "xiaomi",
			"price":      int64(100),
			"is_active":  true,
		},
		{
			"entity_key": "1234",
			"unix_time":  int64(20),
			"model":      "galaxy",
			"price":      int64(300),
			"is_active":  true,
		},
		{
			"entity_key": "1235",
			"unix_time":  int64(5),
			"model":      "apple",
			"price":      int64(200),
			"is_active":  false,
		},
		{
			"entity_key": "1235",
			"unix_time":  int64(15),
			"model":      "oneplus",
			"price":      int64(240),
			"is_active":  false,
		},
	}
	if oneGroup {
		header = header[:len(header)-1]
	}

	data := make(chan []interface{})
	go func() {
		defer close(data)
		for _, item := range values {
			record := make([]interface{}, 0, len(header))
			for _, h := range header {
				record = append(record, item[h])
			}
			data <- record
		}
	}()
	return &types.JoinResult{
		Header: header,
		Data:   data,
	}
}

func extractValues(stream <-chan []interface{}) [][]interface{} {
	values := make([][]interface{}, 0)
	for item := range stream {
		values = append(values, item)
	}
	return values
}
