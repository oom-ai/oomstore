package test_impl

import (
	"encoding/csv"
	"sort"
	"strings"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entity := types.Entity{
		Name:   "device",
		Length: 16,
	}
	dataTable := "data_1_1"

	opt := offline.ImportOpt{
		Entity:        &entity,
		DataTableName: dataTable,
		Features: []*types.Feature{
			{
				Name:        "model",
				DBValueType: "invalid-db-value-type",
			},
			{
				Name:        "price",
				DBValueType: "int",
			},
		},
		Header: []string{"device", "model", "price"},
		CsvReader: csv.NewReader(strings.NewReader(`
1234,xiaomi,1899
1235,apple,4999
1236,huawei,5999
1237,oneplus,3999
`)),
	}

	t.Run("invalid db value type", func(t *testing.T) {
		_, err := store.Import(ctx, opt)
		assert.Error(t, err)
	})

	t.Run("normal import call", func(t *testing.T) {
		revision := int64(1234)
		opt.Features[0].DBValueType = "varchar(32)"
		opt.Revision = &revision
		_, err := store.Import(ctx, opt)
		assert.NoError(t, err)

		stream, err := store.Export(ctx, offline.ExportOpt{
			DataTable:    dataTable,
			EntityName:   entity.Name,
			FeatureNames: []string{"model", "price"},
		})
		assert.NoError(t, err)

		records := make([][]interface{}, 0)
		for ele := range stream {
			records = append(records, ele.Record)
		}
		sort.Slice(records, func(i, j int) bool {
			return cast.ToString(records[i][0]) < cast.ToString(records[j][0])
		})
		assert.Equal(t, records, [][]interface{}{
			{"1234", "xiaomi", int64(1899)},
			{"1235", "apple", int64(4999)},
			{"1236", "huawei", int64(5999)},
			{"1237", "oneplus", int64(3999)},
		})
	})
}
