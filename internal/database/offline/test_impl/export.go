package test_impl

import (
	"bufio"
	"fmt"
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

	snapshotTable := "offline_snapshot_1_1"
	cdcTable := "offline_cdc_1_1"
	unixMilli := &types.Feature{
		Name:      "unix_milli",
		ValueType: types.Int64,
	}
	features := []*types.Feature{
		{
			Name:      "model",
			ValueType: types.String,
		},
		{
			Name:      "price",
			ValueType: types.Int64,
		},
	}
	buildTestSnapshotTable(ctx, t, store, features, 1, snapshotTable, &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,100
1235,apple,200
1236,huawei,300
1237,oneplus,240
`)),
		Delimiter: ",",
	})
	buildTestSnapshotTable(ctx, t, store, append(features, unixMilli), 1, cdcTable, &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi-1,120,2
1235,apple-2,115,5
1234,xiaomi-1,130,10
1238,galaxy,100,11
1239,galaxy,90,12
`)),
		Delimiter: ",",
	})

	testCases := []struct {
		description   string
		opt           offline.ExportOpt
		expected      [][]interface{}
		expectedError error
	}{
		{
			description: "no features",
			opt: offline.ExportOpt{
				SnapshotTable: snapshotTable,
				EntityName:    "device",
				Features:      types.FeatureList{},
			},
			expected: [][]interface{}{{"1234"}, {"1235"}, {"1236"}, {"1237"}},
		},
		{
			description: "invalid option",
			opt: offline.ExportOpt{
				SnapshotTable: snapshotTable,
				CdcTable:      &cdcTable,
				EntityName:    "device",
				Features:      features,
			},
			expectedError: fmt.Errorf("invalid option %+v", offline.ExportOpt{
				SnapshotTable: snapshotTable,
				CdcTable:      &cdcTable,
				EntityName:    "device",
				Features:      features,
			}),
		},
		{
			description: "batch features",
			opt: offline.ExportOpt{
				SnapshotTable: snapshotTable,
				EntityName:    "device",
				Features:      features,
			},
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
		{
			description: "streaming features",
			opt: offline.ExportOpt{
				SnapshotTable: snapshotTable,
				CdcTable:      &cdcTable,
				UnixMilli:     int64Ptr(11),
				EntityName:    "device",
				Features:      features,
			},
			expected: [][]interface{}{{"1234", "xiaomi-1", int64(130)}, {"1235", "apple-2", int64(115)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}, {"1238", "galaxy", int64(100)}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, errs := store.Export(ctx, tc.opt)
			values := make([][]interface{}, 0)
			for row := range actual {
				values = append(values, row)
			}
			if tc.expectedError != nil {
				assert.EqualError(t, <-errs, tc.expectedError.Error())
			} else {
				assert.ElementsMatch(t, tc.expected, values)
				assert.NoError(t, <-errs)
			}
		})
	}
}

func int64Ptr(i int64) *int64 {
	return &i
}
