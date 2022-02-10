package test_impl

import (
	"bufio"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestImport(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entity := types.Entity{Name: "device"}
	group := &types.Group{
		ID:       1,
		Name:     "device",
		Category: types.CategoryBatch,
	}
	snapshotTable := "offline_1_1"

	opt := offline.ImportOpt{
		EntityName:        entity.Name,
		SnapshotTableName: snapshotTable,
		Header:            []string{"device", "model", "price"},
		Source: &offline.CSVSource{
			Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,1899
1235,apple,4999
1236,huawei,5999
1237,oneplus,3999
`)),
			Delimiter: ',',
		},
	}

	t.Run("normal import call", func(t *testing.T) {
		revision := int64(1234)
		opt.Features = []*types.Feature{
			{
				Name:      "price",
				ValueType: types.Int64,
			},
			{
				Name:      "model",
				ValueType: types.String,
			},
		}
		opt.Revision = &revision
		_, err := store.Import(ctx, opt)
		assert.NoError(t, err)

		result, err := store.Export(ctx, offline.ExportOpt{
			SnapshotTables: map[int]string{1: snapshotTable},
			EntityName:     entity.Name,
			Features: map[int]types.FeatureList{
				1: []*types.Feature{
					{Name: "model", ValueType: types.String, Group: group},
					{Name: "price", ValueType: types.Int64, Group: group},
				}},
		})
		records := make([][]interface{}, 0)
		for row := range result.Data {
			assert.NoError(t, row.Error)
			records = append(records, row.Record)
		}
		assert.NoError(t, err)
		sort.Slice(records, func(i, j int) bool {
			return cast.ToString(records[i][0]) < cast.ToString(records[j][0])
		})
		assert.Equal(t, [][]interface{}{
			{"1234", "xiaomi", int64(1899)},
			{"1235", "apple", int64(4999)},
			{"1236", "huawei", int64(5999)},
			{"1237", "oneplus", int64(3999)},
		}, records)
	})
}
