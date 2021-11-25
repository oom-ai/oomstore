package test_impl

import (
	"encoding/csv"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestExport(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	dataTable := "offline_1_1"
	features := []*types.Feature{
		{
			Name:        "model",
			DBValueType: "VARCHAR(32)"},
		{
			Name:        "price",
			DBValueType: "INT",
		},
	}
	buildTestDataTable(ctx, t, store, features, dataTable, csv.NewReader(strings.NewReader(`
1234,xiaomi,100
1235,apple,200
1236,huawei,300
1237,oneplus,240
`)))

	testCases := []struct {
		description string
		opt         offline.ExportOpt
		expected    [][]interface{}
	}{
		{
			description: "no features",
			opt: offline.ExportOpt{
				DataTable:    dataTable,
				EntityName:   "device",
				FeatureNames: []string{},
			},
			expected: [][]interface{}{{"1234"}, {"1235"}, {"1236"}, {"1237"}},
		},
		{
			description: "valid features and valid entity rows",
			opt: offline.ExportOpt{
				DataTable:    dataTable,
				EntityName:   "device",
				FeatureNames: []string{"model", "price"},
			},
			expected: [][]interface{}{{"1234", "xiaomi", int64(100)}, {"1235", "apple", int64(200)}, {"1236", "huawei", int64(300)}, {"1237", "oneplus", int64(240)}},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.description, func(t *testing.T) {
			actual, errs := store.Export(ctx, tc.opt)
			values := make([][]interface{}, 0)
			for row := range actual {
				values = append(values, row)
			}
			assert.Equal(t, tc.expected, values)
			assert.NoError(t, <-errs)
		})
	}
}
