package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

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
		EntityName: entity.Name,
		GroupID:    group.ID,
	}))

	assert.NoError(t, store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
		EntityName: entity.Name,
		GroupID:    group.ID,
		Feature:    feature1,
	}))

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1},
		FeatureValues: []interface{}{"post1"},
	}))
	rs, err := store.Get(ctx, online.GetOpt{
		EntityKey:  "user1",
		Group:      *group,
		Features:   types.FeatureList{feature1},
		RevisionID: nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName(): "post1",
	}, rs)

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1},
		FeatureValues: []interface{}{"post2"},
	}))
	rs, err = store.Get(ctx, online.GetOpt{
		EntityKey:  "user1",
		Group:      *group,
		Features:   types.FeatureList{feature1},
		RevisionID: nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName(): "post2",
	}, rs)

	assert.NoError(t, store.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
		EntityName: entity.Name,
		GroupID:    group.ID,
		Feature:    feature2,
	}))

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1, feature2},
		FeatureValues: []interface{}{"post1", "post2"},
	}))
	rs, err = store.Get(ctx, online.GetOpt{
		EntityKey:  "user1",
		Group:      *group,
		Features:   types.FeatureList{feature1, feature2},
		RevisionID: nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName(): "post1",
		feature2.FullName(): "post2",
	}, rs)
}
