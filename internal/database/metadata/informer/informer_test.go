package informer_test

import (
	"context"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata/informer"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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
	entities := types.EntityList{&entity}
	groups := types.GroupList{&group}
	return informer.NewCache(entities, nil, groups, nil)
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

	entity, err := informer.CacheGetEntity(ctx, 1)
	require.NoError(t, err)

	// changing this entity should not change the internal state of the informer
	entity.Length = 20

	entity, err = informer.CacheGetEntity(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, 1, entity.ID)
	require.Equal(t, 10, entity.Length)
	require.Equal(t, "entity", entity.Name)
}

func TestInformerDeepCopyList(t *testing.T) {
	ctx, informer := prepareInformer(t)
	defer informer.Close()

	entityList := informer.CacheListEntity(ctx)
	require.Equal(t, 1, len(entityList))

	// changing this entity should not change the internal state of the informer
	entityList[0].Length = 20

	entity, err := informer.CacheGetEntity(ctx, 1)
	require.NoError(t, err)
	require.Equal(t, 1, entity.ID)
	require.Equal(t, 10, entity.Length)
	require.Equal(t, "entity", entity.Name)
}
