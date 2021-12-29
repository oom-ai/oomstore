package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

type StreamData struct {
	entity   types.Entity
	groups   types.GroupList
	features types.FeatureList
}

var simpleStreamData StreamData

func init() {
	entity := types.Entity{
		ID:     1,
		Name:   "user",
		Length: 16,
	}

	group1 := &types.Group{
		ID:       1,
		Name:     "user_clicks",
		Category: types.CategoryStream,
		EntityID: 1,
		Entity:   &entity,
	}

	group2 := &types.Group{
		ID:       2,
		Name:     "user_reads",
		Category: types.CategoryStream,
		EntityID: 1,
		Entity:   &entity,
	}

	simpleStreamData = StreamData{
		entity: entity,
		groups: types.GroupList{group1, group2},
		features: types.FeatureList{
			&types.Feature{
				ID:        1,
				Name:      "last_5_click_post",
				FullName:  "user_clicks.last_5_click_post",
				ValueType: types.String,
				GroupID:   1,
				Group:     group1,
			},
			&types.Feature{
				ID:        2,
				Name:      "last_10_click_post",
				FullName:  "user_clicks.last_10_click_post",
				ValueType: types.String,
				GroupID:   1,
				Group:     group1,
			},
			&types.Feature{
				ID:        3,
				Name:      "recent_5_read_post",
				FullName:  "user_reads.recent_5_read_post",
				ValueType: types.String,
				GroupID:   2,
				Group:     group2,
			},
		},
	}
}

func TestPrepareStreamTable(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	for _, group := range simpleStreamData.groups {
		t.Run("create stream table", func(t *testing.T) {
			err := store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
				Entity:  group.Entity,
				GroupID: group.ID,
			})
			assert.NoError(t, err, "create stream table failed: %v", err)
		})
	}

	for _, feature := range simpleStreamData.features {
		t.Run("stream table add column", func(t *testing.T) {
			err := store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
				Entity:  feature.Entity(),
				GroupID: feature.GroupID,
				Feature: feature,
			})
			assert.NoError(t, err, "stream table add column failed: %v", err)
		})
	}
}