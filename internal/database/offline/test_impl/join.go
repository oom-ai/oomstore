package test_impl

import (
	"bufio"
	"context"
	"math"
	"sort"
	"strings"
	"testing"
	"time"

	"github.com/spf13/cast"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestJoin(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entity := &types.Entity{Name: "device"}
	unixMilli := &types.Feature{
		Name:      "unix_milli",
		ValueType: types.Int64,
	}
	oneGroupFeatures, oneGroupFeatureMap, oneGroupGroupNames := prepareFeatures(true)
	twoGroupFeatures, twoGroupFeatureMap, twoGroupGroupNames := prepareFeatures(false)

	buildTestSnapshotTable(ctx, t, store, oneGroupFeatures, 1, "offline_snapshot_1_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,100
1235,apple,200
`)),
		Delimiter: ',',
	})
	buildTestSnapshotTable(ctx, t, store, oneGroupFeatures, 2, "offline_snapshot_1_2", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,galaxy,300
1235,oneplus,240
`)),
		Delimiter: ',',
	})
	buildTestSnapshotTable(ctx, t, store, twoGroupFeatures[2:], 1, "offline_snapshot_2_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,1
1235,0
`)),
		Delimiter: ',',
	})

	buildTestSnapshotTable(ctx, t, store, append(oneGroupFeatures, unixMilli), 1, "offline_cdc_1_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi-1,120,2
1235,apple-2,115,14
1234,xiaomi-1,130,10
`)),
		Delimiter: ',',
	})
	buildTestSnapshotTable(ctx, t, store, append(oneGroupFeatures, unixMilli), 2, "offline_cdc_1_2", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,galaxy-1,320,18
`)),
		Delimiter: ',',
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
			expected: prepareEmptyResult(),
		},
		{
			description: "no entity rows",
			opt: offline.JoinOpt{
				EntityName: entity.Name,
				GroupNames: oneGroupGroupNames,
				FeatureMap: oneGroupFeatureMap,
				EntityRows: prepareEntityRows(true, false),
			},
			expected: prepareEmptyResult(),
		},
		{
			description: "one batch feature group",
			opt: offline.JoinOpt{
				EntityName:       entity.Name,
				EntityRows:       prepareEntityRows(false, false),
				GroupNames:       oneGroupGroupNames,
				FeatureMap:       oneGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(true, false),
			},
			expected: prepareResult(true, false, prepareBatchResultValues()),
		},
		{
			description: "two batch feature groups",
			opt: offline.JoinOpt{
				EntityName:       entity.Name,
				EntityRows:       prepareEntityRows(false, false),
				GroupNames:       twoGroupGroupNames,
				FeatureMap:       twoGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(false, false),
			},
			expected: prepareResult(false, false, prepareBatchResultValues()),
		},
		{
			description: "two batch feature groups, with extra values",
			opt: offline.JoinOpt{
				EntityName:       entity.Name,
				EntityRows:       prepareEntityRows(false, true),
				GroupNames:       twoGroupGroupNames,
				FeatureMap:       twoGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(false, false),
				ValueNames:       []string{"value_1", "value_2"},
			},
			expected: prepareResult(false, true, prepareBatchResultValues()),
		},
		{
			description: "one streaming feature group, one batch group",
			opt: offline.JoinOpt{
				EntityName:       entity.Name,
				EntityRows:       prepareEntityRows(false, true),
				GroupNames:       twoGroupGroupNames,
				FeatureMap:       twoGroupFeatureMap,
				RevisionRangeMap: prepareRevisionRanges(false, true),
				ValueNames:       []string{"value_1", "value_2"},
			},
			expected: prepareResult(false, true, prepareStreamingResultValues()),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, err := store.Join(context.Background(), tc.opt)
			require.NoError(t, err)

			assert.ElementsMatch(t, tc.expected.Header, actual.Header)
			expectedValues := extractValues(t, tc.expected.Data)
			actualValues := extractValues(t, actual.Data)
			for i := range expectedValues {
				assert.ElementsMatch(t, expectedValues[i], actualValues[i])
			}

			tempTables, err := store.GetTemporaryTables(ctx, time.Now().UnixMilli())
			assert.NoError(t, err)
			assert.Equal(t, 0, len(tempTables))
		})
	}
}

func prepareFeatures(oneGroup bool) (types.FeatureList, map[string]types.FeatureList, []string) {
	groupBasic := &types.Group{
		ID:       1,
		Name:     "device_basic",
		Entity:   &types.Entity{Name: "device"},
		Category: types.CategoryBatch,
	}
	groupAdvanced := &types.Group{
		ID:       2,
		Name:     "device_advanced",
		Entity:   &types.Entity{Name: "device"},
		Category: types.CategoryBatch,
	}
	price := &types.Feature{
		Name:      "price",
		ValueType: types.Int64,
		GroupID:   1,
		Group:     groupBasic,
	}
	model := &types.Feature{
		Name:      "model",
		ValueType: types.String,
		GroupID:   1,
		Group:     groupBasic,
	}
	isActive := &types.Feature{
		Name:      "is_active",
		ValueType: types.Bool,
		GroupID:   2,
		Group:     groupAdvanced,
	}

	if oneGroup {
		features := types.FeatureList{model, price}
		featureMap := map[string]types.FeatureList{
			groupBasic.Name: {model, price},
		}
		return features, featureMap, []string{groupBasic.Name}
	} else {
		features := types.FeatureList{model, price, isActive}
		featureMap := map[string]types.FeatureList{
			groupBasic.Name:     {model, price},
			isActive.Group.Name: {isActive},
		}
		return features, featureMap, []string{groupBasic.Name, groupAdvanced.Name}
	}
}

func prepareRevisionRanges(oneGroup bool, stream bool) map[string][]*offline.RevisionRange {
	basic := []*offline.RevisionRange{
		{
			MinRevision:   5,
			MaxRevision:   15,
			SnapshotTable: "offline_snapshot_1_1",
		},
		{
			MinRevision:   15,
			MaxRevision:   25,
			SnapshotTable: "offline_snapshot_1_2",
		},
	}
	advanced := []*offline.RevisionRange{
		{
			MinRevision:   1,
			MaxRevision:   math.MaxInt64,
			SnapshotTable: "offline_snapshot_2_1",
		},
	}
	if stream {
		basic[0].CdcTable = "offline_cdc_1_1"
		basic[1].CdcTable = "offline_cdc_1_2"
	}
	if oneGroup {
		return map[string][]*offline.RevisionRange{
			"device_basic": basic,
		}
	}

	return map[string][]*offline.RevisionRange{
		"device_basic":    basic,
		"device_advanced": advanced,
	}
}

func prepareEntityRows(isEmpty bool, withValue bool) <-chan types.EntityRow {
	entityRows := make(chan types.EntityRow)
	rows := []types.EntityRow{
		{
			EntityKey: "1234",
			UnixMilli: 2,
			Values:    []string{"1", "2"},
		},
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

func prepareResult(oneGroup bool, withValue bool, values []map[string]interface{}) *types.JoinResult {
	header := []string{"entity_key", "unix_milli", "device_basic.model", "device_basic.price", "device_advanced.is_active"}
	if withValue {
		header = []string{"entity_key", "unix_milli", "value_1", "value_2", "device_basic.model", "device_basic.price", "device_advanced.is_active"}
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

	data := make(chan types.JoinRecord)
	go func() {
		defer close(data)
		for _, item := range values {
			record := make([]interface{}, 0, len(header))
			for _, h := range header {
				record = append(record, item[h])
			}
			data <- types.JoinRecord{Record: record}
		}
	}()
	return &types.JoinResult{
		Header: header,
		Data:   data,
	}
}

func prepareEmptyResult() *types.JoinResult {
	data := make(chan types.JoinRecord)
	defer close(data)
	return &types.JoinResult{
		Data: data,
	}
}

func prepareStreamingResultValues() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"entity_key":                "1234",
			"unix_milli":                int64(2),
			"value_1":                   "1",
			"value_2":                   "2",
			"device_basic.model":        nil,
			"device_basic.price":        nil,
			"device_advanced.is_active": true,
		},
		{
			"entity_key":                "1234",
			"unix_milli":                int64(10),
			"value_1":                   "1",
			"value_2":                   "2",
			"device_basic.model":        "xiaomi-1",
			"device_basic.price":        int64(130),
			"device_advanced.is_active": true,
		},
		{
			"entity_key":                "1234",
			"unix_milli":                int64(20),
			"value_1":                   "3",
			"value_2":                   "4",
			"device_basic.model":        "galaxy-1",
			"device_basic.price":        int64(320),
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
			"device_basic.model":        "apple-2",
			"device_basic.price":        int64(115),
			"device_advanced.is_active": false,
		},
	}
}

func prepareBatchResultValues() []map[string]interface{} {
	return []map[string]interface{}{
		{
			"entity_key":                "1234",
			"unix_milli":                int64(2),
			"value_1":                   "1",
			"value_2":                   "2",
			"device_basic.model":        nil,
			"device_basic.price":        nil,
			"device_advanced.is_active": true,
		},
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
}

func extractValues(t *testing.T, stream <-chan types.JoinRecord) [][]interface{} {
	values := make([][]interface{}, 0)
	for item := range stream {
		assert.NoError(t, item.Error)
		values = append(values, item.Record)
	}
	sort.Slice(values, func(i, j int) bool {
		if cast.ToString(values[i][0]) == cast.ToString(values[j][0]) {
			return cast.ToInt(values[i][1]) < cast.ToInt(values[j][1])
		}
		return cast.ToString(values[i][0]) < cast.ToString(values[j][0])
	})
	return values
}
