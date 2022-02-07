package test_impl

import (
	"testing"

	"github.com/oom-ai/oomstore/internal/database/dbutil"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestPush(t *testing.T, prepareStore PrepareStoreFn, destoryStore DestroyStoreFn) {
	t.Cleanup(destoryStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entity := SampleStream.Entity
	group := SampleStream.Group
	feature1 := SampleStream.Features[0]
	feature2 := SampleStream.Features[1]

	assert.NoError(t, store.CreateTable(ctx, online.CreateTableOpt{
		EntityName: entity.Name,
		TableName:  dbutil.OnlineStreamTableName(group.ID),
		Features:   SampleStream.Features,
	}))

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1},
		FeatureValues: []interface{}{int64(1)},
	}))
	rs, err := store.Get(ctx, online.GetOpt{
		EntityKey: "user1",
		Group:     group,
		Features:  types.FeatureList{feature1},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName(): int64(1),
	}, rs)

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1},
		FeatureValues: []interface{}{int64(2)},
	}))
	rs, err = store.Get(ctx, online.GetOpt{
		EntityKey: "user1",
		Group:     group,
		Features:  types.FeatureList{feature1},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName(): int64(2),
	}, rs)

	assert.NoError(t, store.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     "user1",
		GroupID:       group.ID,
		Features:      types.FeatureList{feature1, feature2},
		FeatureValues: []interface{}{int64(3), "post2"},
	}))
	rs, err = store.Get(ctx, online.GetOpt{
		EntityKey: "user1",
		Group:     group,
		Features:  types.FeatureList{feature1, feature2},
	})
	assert.NoError(t, err)
	assert.Equal(t, map[string]interface{}{
		feature1.FullName(): int64(3),
		feature2.FullName(): "post2",
	}, rs)
}
