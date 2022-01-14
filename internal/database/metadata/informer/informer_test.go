package informer_test

import (
	"context"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/util"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func sampleCache() *informer.Cache {
	entity := types.Entity{
		ID:   1,
		Name: "entity",
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
	return informer.NewCache(entities, features, groups)
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

	name := util.ComposeFullFeatureName("group", "price")
	features := informer.ListCachedFeature(ctx, &[]string{name})
	require.Equal(t, 1, len(features))
	feature := features[0]

	assert.Equal(t, 1, feature.ID)
	assert.Equal(t, "price", feature.Name)
	assert.Equal(t, 100, feature.GroupID)

	assert.NotNil(t, feature.Group)
	assert.Equal(t, 100, feature.Group.ID)
	assert.Equal(t, 1, feature.Group.Entity.ID)
	assert.Equal(t, "entity", feature.Group.Entity.Name)
}

func TestInformerDeepCopy(t *testing.T) {
	ctx, informer := prepareInformer(t)
	defer informer.Close()

	features := informer.ListCachedFeature(ctx, nil)
	require.Equal(t, 1, len(features))
	// changing this entity should not change the internal state of the informer
	features[0].Name = "new_price"

	features = informer.ListCachedFeature(ctx, nil)
	require.Equal(t, 1, len(features))
	assert.Equal(t, 1, features[0].ID)
	assert.Equal(t, 100, features[0].GroupID)
	assert.Equal(t, "price", features[0].Name)
}
