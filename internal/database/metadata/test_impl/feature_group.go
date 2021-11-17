package test_impl

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func prepareEntity(t *testing.T, ctx context.Context, store metadata.Store, name string) int {
	entityID, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  name,
			Length:      32,
			Description: "description",
		},
	})
	require.NoError(t, err)
	return entityID
}

func TestGetFeatureGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	id, err := store.CreateFeatureGroup(ctx, opt)
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	// get non-exist feature group
	_, err = store.GetFeatureGroup(ctx, 0)
	require.Error(t, err)

	// get existing feature group
	featureGroup, err := store.GetFeatureGroup(ctx, id)
	require.NoError(t, err)
	require.Equal(t, opt.Name, featureGroup.Name)
	require.Equal(t, opt.EntityID, featureGroup.EntityID)
	require.Equal(t, opt.Description, featureGroup.Description)
	require.Equal(t, opt.Category, featureGroup.Category)
}

func TestListFeatureGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	deviceEntityID := prepareEntity(t, ctx, store, "device")
	userEntityID := prepareEntity(t, ctx, store, "user")

	deviceOpt := metadata.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    deviceEntityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	userBaseOpt := metadata.CreateFeatureGroupOpt{
		Name:        "user_baseinfo",
		EntityID:    userEntityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	userBehaviorOpt := metadata.CreateFeatureGroupOpt{
		Name:        "user_behaviorinfo",
		EntityID:    userEntityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	deviceGroupID, err := store.CreateFeatureGroup(ctx, deviceOpt)
	require.NoError(t, err)
	userGroupID, err := store.CreateFeatureGroup(ctx, userBaseOpt)
	require.NoError(t, err)
	_, err = store.CreateFeatureGroup(ctx, userBehaviorOpt)
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	require.Equal(t, 1, len(store.ListFeatureGroup(ctx, &deviceGroupID)))
	require.Equal(t, 2, len(store.ListFeatureGroup(ctx, &userGroupID)))
	require.Equal(t, 3, len(store.ListFeatureGroup(ctx, nil)))
	require.Equal(t, 0, len(store.ListFeatureGroup(ctx, intPtr(0))))
}

func TestCreateFeatureGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}

	// create successfully
	featureGroupID, err := store.CreateFeatureGroup(ctx, opt)
	require.NotEqual(t, 0, featureGroupID)
	require.NoError(t, err)

	// cannot create feature group with same name
	_, err = store.CreateFeatureGroup(ctx, opt)
	require.Equal(t, fmt.Errorf("feature group device_baseinfo already exists"), err)

	// cannot create feature group with invalid category
	opt.Category = "invalid-category"
	_, err = store.CreateFeatureGroup(ctx, opt)
	require.NotNil(t, err)
}

func TestUpdateFeatureGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	featureGroupID, err := store.CreateFeatureGroup(ctx, opt)
	require.NoError(t, err)

	// update non-exist feature group
	require.NotNil(t, store.UpdateFeatureGroup(ctx, metadata.UpdateFeatureGroupOpt{
		GroupID: 0,
	}))

	// update existing feature group
	description := "new description"
	require.Nil(t, store.UpdateFeatureGroup(ctx, metadata.UpdateFeatureGroupOpt{
		GroupID:        featureGroupID,
		NewDescription: &description,
	}))
}
