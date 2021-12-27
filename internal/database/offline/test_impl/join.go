package test_impl

import (
	"bufio"
	"context"
	"math"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestJoin(t *testing.T, prepareStore PrepareStoreFn, destoryStore DestroyStoreFn) {
	t.Cleanup(destoryStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entity := &types.Entity{
		Name:   "device",
		Length: 10,
	}
	oneGroupFeatures, oneGroupFeatureMap := prepareFeatures(true)
	twoGroupFeatures, twoGroupFeatureMap := prepareFeatures(false)

	buildTestDataTable(ctx, t, store, oneGroupFeatures, "offline_1_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,100
1235,apple,200
`)),
		Delimiter: ",",
	})
	buildTestDataTable(ctx, t, store, oneGroupFeatures, "offline_1_2", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,galaxy,300
1235,oneplus,240
`)),
		Delimiter: ",",
	})
	buildTestDataTable(ctx, t, store, twoGroupFeatures[2:], "offline_2_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,1
1235,0
`)),
		Delimiter: ",",
	})

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
				EntityRows: prepareEntityRows(true, false),
			},
			expected: nil,
		},
		{
			description: "one feature group",
			opt: offline.JoinOpt{
				Entity:           *entity,
				EntityRows:       prepareEntityRows(false, false),
				FeatureMap:       oneGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(true),
			},
			expected: prepareResult(true, false),
		},
		{
			description: "two feature groups",
			opt: offline.JoinOpt{
				Entity:           *entity,
				EntityRows:       prepareEntityRows(false, false),
				FeatureMap:       twoGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(false),
			},
			expected: prepareResult(false, false),
		},
		{
			description: "two feature groups, with extra values",
			opt: offline.JoinOpt{
				Entity:           *entity,
				EntityRows:       prepareEntityRows(false, true),
				FeatureMap:       twoGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(false),
				ValueNames:       []string{"value_1", "value_2"},
			},
			expected: prepareResult(false, true),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := store.Join(context.Background(), tc.opt)
			require.NoError(t, err)
			if tc.expected == nil {
				assert.Equal(t, tc.expected, actual)
			} else {
				expectedValues := extractValues(tc.expected.Data)
				actualValues := extractValues(actual.Data)
				assert.ElementsMatch(t, tc.expected.Header, actual.Header)
				assert.ObjectsAreEqual(expectedValues, actualValues)
			}
		})
	}
}

func prepareFeatures(oneGroup bool) (types.FeatureList, map[string]types.FeatureList) {
	price := &types.Feature{
		Name:      "price",
		FullName:  "device_basic.price",
		ValueType: types.Int64,
		GroupID:   1,
	}
	model := &types.Feature{
		Name:      "model",
		FullName:  "device_basic.model",
		ValueType: types.String,
		GroupID:   1,
	}
	isActive := &types.Feature{
		Name:      "is_active",
		FullName:  "device_advanced.is_active",
		ValueType: types.Bool,
		GroupID:   2,
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

func prepareEntityRows(isEmpty bool, withValue bool) <-chan types.EntityRow {
	entityRows := make(chan types.EntityRow)
	rows := []types.EntityRow{
		{
			EntityKey: "1234",
			UnixMilli: 10,
			Values:    []string{"1", "2"},
		},
		{
			EntityKey: "1234",
			UnixMilli: 20,
			Values:    []string{"3", "4"},
		},
		{
			EntityKey: "1235",
			UnixMilli: 5,
			Values:    []string{"5", "6"},
		},
		{
			EntityKey: "1235",
			UnixMilli: 14,
			Values:    []string{"7", "8"},
		},
	}
	go func() {
		defer close(entityRows)
		if isEmpty {
			return
		}
		for _, row := range rows {
			if !withValue {
				row.Values = nil
			}
			entityRows <- row
		}
	}()
	return entityRows
}

func prepareResult(oneGroup bool, withValue bool) *types.JoinResult {
	header := []string{"entity_key", "unix_milli", "device_basic.model", "device_basic.price", "device_advanced.is_active"}
	if withValue {
		header = []string{"entity_key", "unix_milli", "value_1", "value_2", "device_basic.model", "device_basic.price", "device_advanced.is_active"}
	}
	values := []map[string]interface{}{
		{
			"entity_key":                "1234",
			"unix_milli":                int64(10),
			"value_1":                   "1",
			"value_2":                   "2",
			"device_basic.model":        "xiaomi",
			"device_basic.price":        int64(100),
			"device_advanced.is_active": true,
		},
		{
			"entity_key":                "1234",
			"unix_milli":                int64(20),
			"value_1":                   "3",
			"value_2":                   "4",
			"device_basic.model":        "galaxy",
			"device_basic.price":        int64(300),
			"device_advanced.is_active": true,
		},
		{
			"entity_key":                "1235",
			"unix_milli":                int64(5),
			"value_1":                   "5",
			"value_2":                   "6",
			"device_basic.model":        "apple",
			"device_basic.price":        int64(200),
			"device_advanced.is_active": false,
		},
		{
			"entity_key":                "1235",
			"unix_milli":                int64(14),
			"value_1":                   "7",
			"value_2":                   "8",
			"device_basic.model":        "apple",
			"device_basic.price":        int64(200),
			"device_advanced.is_active": false,
		},
	}
	if oneGroup {
		header = header[:len(header)-1]
	}
	if !withValue {
		for _, m := range values {
			delete(m, "value_1")
			delete(m, "value_2")
		}
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
