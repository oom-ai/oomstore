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

	group := &types.Group{
		ID:     1,
		Entity: &types.Entity{Name: "device"},
	}
	unixMilli := &types.Feature{
		Name:      "unix_milli",
		ValueType: types.Int64,
	}
	features, _ := prepareFeatures(true)

	buildTestSnapshotTable(ctx, t, store, features, 1, "offline_stream_snapshot_1_1", &offline.CSVSource{
		Reader: bufio.NewReader(strings.NewReader(`1234,xiaomi,100
1235,apple,200
1236,oneplus,155
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

	actual, errs := store.Export(ctx, offline.ExportOpt{
		SnapshotTable: "offline_stream_snapshot_1_2",
		EntityName:    "device",
		Features:      features,
	})
	values := make([][]interface{}, 0)
	for row := range actual {
		values = append(values, row)
	}
	sort.Slice(values, func(i, j int) bool {
		return cast.ToInt64(values[i][0]) < cast.ToInt64(values[j][0])
	})
	expected := [][]interface{}{
		{"1234", "xiaomi-1", int64(130)},
		{"1235", "apple-1", int64(113)},
		{"1236", "oneplus", int64(155)},
		{"1237", "pixel", int64(200)},
	}
	assert.Equal(t, expected, values)
	assert.NoError(t, <-errs)

}
