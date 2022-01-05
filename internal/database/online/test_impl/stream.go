package test_impl

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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

func TestPush(t *testing.T, prepareStore PrepareStoreFn, destoryStore DestroyStoreFn) {
	t.Cleanup(destoryStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	var (
		entity   = simpleStreamData.entity
		group    = simpleStreamData.groups[0]
		feature1 = simpleStreamData.features[0]
		feature2 = simpleStreamData.features[1]
	)

	assert.NoError(t, store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
		Entity:  &entity,
		GroupID: group.ID,
	}))

	assert.NoError(t, store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
		Entity:  &entity,
		GroupID: group.ID,
		Feature: feature1,
	}))

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		Entity:        &entity,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1},
		FeatureValues: []interface{}{"post1"},
	}))
	rs, err := store.Get(ctx, online.GetOpt{
		Entity:     &entity,
		EntityKey:  "user1",
		RevisionID: nil,
		Group:      group,
		Features:   types.FeatureList{feature1},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName: "post1",
	}, rs)

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		Entity:        &entity,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1},
		FeatureValues: []interface{}{"post2"},
	}))
	rs, err = store.Get(ctx, online.GetOpt{
		Entity:     &entity,
		EntityKey:  "user1",
		RevisionID: nil,
		Group:      group,
		Features:   types.FeatureList{feature1},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName: "post2",
	}, rs)

	assert.NoError(t, store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
		Entity:  &entity,
		GroupID: group.ID,
		Feature: feature2,
	}))

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		Entity:        &entity,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1, feature2},
		FeatureValues: []interface{}{"post1", "post2"},
	}))
	rs, err = store.Get(ctx, online.GetOpt{
		Entity:     &entity,
		EntityKey:  "user1",
		RevisionID: nil,
		Group:      group,
		Features:   types.FeatureList{feature1, feature2},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName: "post1",
		feature2.FullName: "post2",
	}, rs)
}
