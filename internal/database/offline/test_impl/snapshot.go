package test_impl

import (
	"bufio"
	"sort"
	"strings"
	"testing"

	"github.com/spf13/cast"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func TestSnapshot(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	features, _ := prepareFeatures(true)
	group := features[0].Group
	unixMilli := &types.Feature{
		Name:      "unix_milli",
		ValueType: types.Int64,
		Group:     group,
	}

	buildTestSnapshotTable(ctx, t, store, append(features, unixMilli), 1, "offline_stream_snapshot_1_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,100,1
1235,apple,200,5
1236,oneplus,155,8
`)),
		Delimiter: ",",
	})

	buildTestSnapshotTable(ctx, t, store, append(features, unixMilli), 2, "offline_stream_cdc_1_2", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi-1,120,2
1235,apple-2,115,14
1234,xiaomi-1,130,10
1237,pixel,200,11
1235,apple-1,113,15
`)),
		Delimiter: ",",
	})

	err := store.Snapshot(ctx, offline.SnapshotOpt{
		Group:        group,
		Features:     features,
		Revision:     2,
		PrevRevision: 1,
	})
	assert.NoError(t, err)

	result, err := store.Export(ctx, offline.ExportOpt{
		SnapshotTables: map[int]string{1: "offline_stream_snapshot_1_2"},
		EntityName:     "device",
		Features:       map[int]types.FeatureList{1: append(features, unixMilli)},
	})
	values := make([][]interface{}, 0)
	for row := range result.Data {
		values = append(values, row)
	}
	sort.Slice(values, func(i, j int) bool {
		return cast.ToInt64(values[i][0]) < cast.ToInt64(values[j][0])
	})
	expected := [][]interface{}{
		{"1234", "xiaomi-1", int64(130), int64(10)},
		{"1235", "apple-1", int64(113), int64(15)},
		{"1236", "oneplus", int64(155), int64(8)},
		{"1237", "pixel", int64(200), int64(11)},
	}
	assert.Equal(t, expected, values)
	assert.NoError(t, result.CheckStreamError())
	assert.NoError(t, err)
}
