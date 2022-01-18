package test_impl

import (
	"bufio"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cast"
	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestImportStream(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entity := types.Entity{Name: "device"}
	group := &types.Group{
		ID:       1,
		Name:     "device_basic",
		Category: types.CategoryStream,
		Entity:   &entity,
	}
	features := []*types.Feature{
		{
			Name:      "price",
			ValueType: types.Int64,
			Group:     group,
		},
		{
			Name:      "model",
			ValueType: types.String,
			Group:     group,
		},
	}

	opt := offline.ImportStreamOpt{
		Entity: &entity,
		Header: []string{"device", "unix_milli", "model", "price"},
		Source: &offline.CSVSource{
			Reader: bufio.NewReader(strings.NewReader(`1234,1,xiaomi,1899
1235,5,apple,4999
1236,8,huawei,5999
1237,10,oneplus,3999
1235,15,apple-1,5999
1238,17,galaxy,6999
`)),
			Delimiter: ",",
		},
	}

	t.Run("import stream features, succeed", func(t *testing.T) {
		opt.Features = features
		timeRange, err := store.ImportStream(ctx, opt)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), timeRange.MinUnixMilli)
		assert.Equal(t, int64(17), timeRange.MaxUnixMilli)

		result, err := store.Export(ctx, offline.ExportOpt{
			SnapshotTables: map[int]string{1: dbutil.OfflineStreamSnapshotTableName(group.ID, 1)},
			CdcTables:      map[int]string{1: dbutil.OfflineStreamCdcTableName(group.ID, 1)},
			EntityName:     entity.Name,
			UnixMilli:      15,
			Features:       map[int]types.FeatureList{1: features},
		})
		assert.NoError(t, err)
		records := make([][]interface{}, 0)
		for row := range result.Data {
			records = append(records, row)
		}
		assert.NoError(t, err)
		assert.NoError(t, result.CheckStreamError())
		sort.Slice(records, func(i, j int) bool {
			return cast.ToString(records[i][0]) < cast.ToString(records[j][0])
		})
		assert.Equal(t, [][]interface{}{
			{"1234", int64(1899), "xiaomi"},
			{"1235", int64(5999), "apple-1"},
			{"1236", int64(5999), "huawei"},
			{"1237", int64(3999), "oneplus"},
		}, records)
	})
}
