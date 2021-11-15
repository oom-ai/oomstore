package test

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func TestCreateEntity(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	opt := metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}
	_, err := store.CreateEntity(ctx, opt)
	require.NoError(t, err)

	_, err = store.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.Equal(t, err, fmt.Errorf("entity device already exists"))
}

func TestGetEntity(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	opt := metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}

	id, err := store.CreateEntity(ctx, opt)
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	entity, err := store.GetEntity(ctx, id)
	require.NoError(t, err)
	require.Equal(t, opt.Name, entity.Name)
	require.Equal(t, opt.Length, entity.Length)
	require.Equal(t, opt.Description, entity.Description)

	_, err = store.GetEntity(ctx, 0)
	require.EqualError(t, err, "feature entity 0 not found")
}

func TestUpdateEntity(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	id, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	require.NoError(t, store.UpdateEntity(ctx, metadata.UpdateEntityOpt{
		EntityID:       id,
		NewDescription: "new description",
	}))

	require.NoError(t, store.Refresh())

	entity, err := store.GetEntity(ctx, id)
	require.NoError(t, err)
	require.Equal(t, entity.Description, "new description")

	require.Error(t, store.UpdateEntity(ctx, metadata.UpdateEntityOpt{
		EntityID:       id + 1,
		NewDescription: "new description",
	}))
}

func TestListEntity(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	require.NoError(t, store.Refresh())

	entitys := store.ListEntity(ctx)
	require.Equal(t, 0, len(entitys))

	_, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	entitys = store.ListEntity(ctx)
	require.Equal(t, 1, len(entitys))
	_, err = store.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "user",
		Length:      16,
		Description: "description",
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	entitys = store.ListEntity(ctx)
	require.Equal(t, 2, len(entitys))
}
