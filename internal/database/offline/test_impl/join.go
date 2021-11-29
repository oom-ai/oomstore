package test_impl

import (
	"context"
	"encoding/csv"
	"math"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestJoin(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entity := &types.Entity{
		Name:   "device",
		Length: 10,
	}
	oneGroupFeatures, oneGroupFeatureMap := prepareFeatures(true)
	twoGroupFeatures, twoGroupFeatureMap := prepareFeatures(false)

	buildTestDataTable(ctx, t, store, oneGroupFeatures, "offline_1_1", csv.NewReader(strings.NewReader(`
1234,xiaomi,100
1235,apple,200
`)))
	buildTestDataTable(ctx, t, store, oneGroupFeatures, "offline_1_2", csv.NewReader(strings.NewReader(`
1234,galaxy,300
1235,oneplus,240
`)))
	buildTestDataTable(ctx, t, store, twoGroupFeatures[2:], "offline_2_1", csv.NewReader(strings.NewReader(`
1234,true
1235,false
`)))

	testCases := []struct {
		description string
		opt         offline.JoinOpt
		expected    *types.JoinResult
	}{
		{
			description: "no features",
			opt: offline.JoinOpt{
				FeatureMap: make(map[string]types.FeatureList),
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
			actual, err := store.Join(context.Background(), tc.opt)
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

func prepareFeatures(oneGroup bool) (types.FeatureList, map[string]types.FeatureList) {
	price := &types.Feature{
		Name:        "price",
		DBValueType: "INT",
		GroupID:     1,
	}
	model := &types.Feature{
		Name:        "model",
		DBValueType: "VARCHAR(32)",
		GroupID:     1,
	}
	isActive := &types.Feature{
		Name:        "is_active",
		DBValueType: "boolean",
		GroupID:     2,
	}

	if oneGroup {
		features := types.FeatureList{model, price}
		featureMap := map[string]types.FeatureList{
			"device_basic": {
				model, price,
			},
		}
		return features, featureMap
	} else {
		features := types.FeatureList{model, price, isActive}
		featureMap := map[string]types.FeatureList{
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
			DataTable:   "offline_1_1",
		},
		{
			MinRevision: 15,
			MaxRevision: 25,
			DataTable:   "offline_1_2",
		},
	}
	advanced := []*metadata.RevisionRange{
		{
			MinRevision: 5,
			MaxRevision: math.MaxInt64,
			DataTable:   "offline_2_1",
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
			UnixMilli: 10,
		}
		entityRows <- types.EntityRow{
			EntityKey: "1234",
			UnixMilli: 20,
		}
		entityRows <- types.EntityRow{
			EntityKey: "1235",
			UnixMilli: 5,
		}
		entityRows <- types.EntityRow{
			EntityKey: "1235",
			UnixMilli: 14,
		}
	}()
	return entityRows
}

func prepareResult(oneGroup bool) *types.JoinResult {
	header := []string{"entity_key", "unix_milli", "model", "price", "is_active"}
	values := []map[string]interface{}{
		{
			"entity_key": "1234",
			"unix_milli": int64(10),
			"model":      "xiaomi",
			"price":      int64(100),
			"is_active":  true,
		},
		{
			"entity_key": "1234",
			"unix_milli": int64(20),
			"model":      "galaxy",
			"price":      int64(300),
			"is_active":  true,
		},
		{
			"entity_key": "1235",
			"unix_milli": int64(5),
			"model":      "apple",
			"price":      int64(200),
			"is_active":  false,
		},
		{
			"entity_key": "1235",
			"unix_milli": int64(15),
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
