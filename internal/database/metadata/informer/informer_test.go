package informer_test

import (
	"context"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleCache() *informer.Cache {
	entity := types.Entity{
		ID:     1,
		Length: 10,
		Name:   "entity",
	}
	group := types.Group{
		ID:       100,
		Name:     "group",
		Category: "batch",
		EntityID: entity.ID,
		Entity:   &entity,
	}
	feature := types.Feature{
		ID:      1,
		Name:    "price",
		GroupID: group.ID,
	}
	entities := types.EntityList{&entity}
	groups := types.GroupList{&group}
	features := types.FeatureList{&feature}
	return informer.NewCache(entities, features, groups, nil)
}

func prepareInformer(t *testing.T) (context.Context, *informer.Informer) {
	ctx := context.Background()

	informer, err := informer.New(time.Second, func() (*informer.Cache, error) {
		return sampleCache(), nil
	})
	require.NoError(t, err)
	return ctx, informer
}

func TestInformer(t *testing.T) {
	ctx, informer := prepareInformer(t)
	defer informer.Close()

	group, err := informer.CacheGetGroup(ctx, 100)
	require.NoError(t, err)

	require.Equal(t, 100, group.ID)
	require.Equal(t, "group", group.Name)
	require.Equal(t, "batch", group.Category)
	require.Equal(t, 1, group.EntityID)

	require.NotNil(t, group.Entity)
	require.Equal(t, 1, group.Entity.ID)
	require.Equal(t, 10, group.Entity.Length)
	require.Equal(t, "entity", group.Entity.Name)
}

func TestInformerDeepCopyGet(t *testing.T) {
	ctx, informer := prepareInformer(t)
	defer informer.Close()

	feature, err := informer.CacheGetFeature(ctx, 1)
	require.NoError(t, err)

	// changing this entity should not change the internal state of the informer
	feature.Name = "new_price"

	feature, err = informer.CacheGetFeature(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, feature.ID)
	assert.Equal(t, 100, feature.GroupID)
	assert.Equal(t, "price", feature.Name)
}

func TestInformerDeepCopyList(t *testing.T) {
	ctx, informer := prepareInformer(t)
	defer informer.Close()

	featureList := informer.CacheListFeature(ctx, metadata.ListFeatureOpt{})
	require.Equal(t, 1, len(featureList))

	// changing this entity should not change the internal state of the informer
	featureList[0].Name = "new_price"

	feature, err := informer.CacheGetFeature(ctx, 1)
	assert.NoError(t, err)
	assert.Equal(t, 1, feature.ID)
	assert.Equal(t, 100, feature.GroupID)
	assert.Equal(t, "price", feature.Name)
}
