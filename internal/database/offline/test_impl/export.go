package test_impl

import (
	"bufio"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)
	ctx, store := prepareStore(t)
	defer store.Close()

	batchSnapshotTable := "offline_batch_snapshot_1_1"
	streamSnapshotTable := "offline_stream_snapshot_2_1"
	streamCdcTable := "offline_stream_cdc_2_1"
	unixMilli := &types.Feature{
		Name:      "unix_milli",
		ValueType: types.Int64,
	}
	batchFeatures, streamFeatures := prepareFeaturesForExport()
	buildTestSnapshotTable(ctx, t, store, batchFeatures, 1, batchSnapshotTable, &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,100
1235,apple,200
1236,huawei,300
1237,oneplus,240
`)),
		Delimiter: ',',
	})
	buildTestSnapshotTable(ctx, t, store, streamFeatures, 1, streamSnapshotTable, &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,1000,true
1235,2040,false
1236,1560,true
1237,4000,false
`)),
		Delimiter: ',',
	})
	buildTestSnapshotTable(ctx, t, store, append(streamFeatures, unixMilli), 1, streamCdcTable, &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,1200,true,2
1235,2050,false,5
1234,1300,false,10
1238,1500,true,11
1239,2700,false,12
`)),
		Delimiter: ',',
	})

	testCases := []struct {
		description string
		opt         offline.ExportOpt
		expected    [][]interface{}
	}{
		{
			description: "one group, batch features",
			opt: offline.ExportOpt{
				SnapshotTables: map[int]string{1: batchSnapshotTable},
				Features:       map[int]types.FeatureList{1: batchFeatures},
				EntityName:     "device",
				UnixMilli:      10,
			},
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "one group, streaming features",
			opt: offline.ExportOpt{
				SnapshotTables: map[int]string{2: streamSnapshotTable},
				CdcTables:      map[int]string{2: streamCdcTable},
				Features:       map[int]types.FeatureList{2: streamFeatures},
				UnixMilli:      11,
				EntityName:     "device",
			},
			expected: [][]interface{}{{"1234", int64(1300), false}, {"1235", int64(2050), false}, {"1236", int64(1560), true}, {"1237", int64(4000), false}, {"1238", int64(1500), true}},
		},
		{
			description: "multiple groups, batch and stream features",
			opt: offline.ExportOpt{
				SnapshotTables: map[int]string{1: batchSnapshotTable, 2: streamSnapshotTable},
				CdcTables:      map[int]string{2: streamCdcTable},
				Features:       map[int]types.FeatureList{1: batchFeatures, 2: streamFeatures},
				EntityName:     "device",
				UnixMilli:      11,
			},
			expected: [][]interface{}{{"1234", "xiaomi", int64(100), int64(1300), false}, {"1235", "apple", int64(200), int64(2050), false}, {"1236", "huawei", int64(300), int64(1560), true}, {"1237", "oneplus", int64(240), int64(4000), false}, {"1238", nil, nil, int64(1500), true}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			result, err := store.Export(ctx, tc.opt)
			values := make([][]interface{}, 0)
			for row := range result.Data {
				assert.NoError(t, row.Error)
				values = append(values, row.Record)
			}
			assert.ElementsMatch(t, tc.expected, values)
			assert.NoError(t, err)
		})
	}
}

func prepareFeaturesForExport() (batchFeatures, streamFeatures types.FeatureList) {
	batchGroup := &types.Group{
		ID:       1,
		Name:     "device",
		Category: types.CategoryBatch,
	}
	streamGroup := &types.Group{
		ID:       2,
		Name:     "account",
		Category: types.CategoryStream,
	}

	batchFeatures = []*types.Feature{
		{
			Name:      "model",
			ValueType: types.String,
			Group:     batchGroup,
		},
		{
			Name:      "price",
			ValueType: types.Int64,
			Group:     batchGroup,
		},
	}
	streamFeatures = []*types.Feature{
		{
			Name:      "last_txn_amount",
			ValueType: types.Int64,
			Group:     streamGroup,
		},
		{
			Name:      "is_vip",
			ValueType: types.Bool,
			Group:     streamGroup,
		},
	}
	return batchFeatures, streamFeatures
}
