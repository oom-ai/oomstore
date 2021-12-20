package test_impl

import (
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func TestCreateEntity(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore()
	defer store.Close()

	opt := metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	}
	_, err := store.CreateEntity(ctx, opt)
	require.NoError(t, err)

	_, err = store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	})
	require.Equal(t, err, fmt.Errorf("entity device already exists"))
}

func TestGetEntity(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore()
	defer store.Close()

	opt := metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	}

	id, err := store.CreateEntity(ctx, opt)
	require.NoError(t, err)

	entity, err := store.GetEntity(ctx, id)
	require.NoError(t, err)
	require.Equal(t, opt.EntityName, entity.Name)
	require.Equal(t, opt.Length, entity.Length)
	require.Equal(t, opt.Description, entity.Description)

	_, err = store.GetEntity(ctx, 0)
	require.EqualError(t, err, "feature entity 0 not found")
}

func TestUpdateEntity(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore()
	defer store.Close()

	id, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	})
	require.NoError(t, err)

	require.NoError(t, store.UpdateEntity(ctx, metadata.UpdateEntityOpt{
		EntityID:       id,
		NewDescription: stringPtr("new description"),
	}))

	require.NoError(t, store.Refresh())

	entity, err := store.GetEntity(ctx, id)
	require.NoError(t, err)
	require.Equal(t, entity.Description, "new description")

	require.Error(t, store.UpdateEntity(ctx, metadata.UpdateEntityOpt{
		EntityID:       id + 1,
		NewDescription: stringPtr("new description"),
	}))
}

func TestListEntity(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore()
	defer store.Close()

	entities, err := store.ListEntity(ctx, nil)
	require.NoError(t, err)
	require.Equal(t, 0, len(entities))

	_, err = store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	})
	require.NoError(t, err)

	entities, err = store.ListEntity(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(entities))

	_, err = store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "user",
			Length:      16,
			Description: "description",
		},
	})
	require.NoError(t, err)

	entities, err = store.ListEntity(ctx, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))

	ids := []int{1, 2}
	entities, err = store.ListEntity(ctx, &ids)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(entities))

	ids = []int{}
	entities, err = store.ListEntity(ctx, &ids)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(entities))
}
