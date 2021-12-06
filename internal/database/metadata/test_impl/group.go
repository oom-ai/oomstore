package test_impl

import (
	"context"
	"fmt"
	"testing"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
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

func TestGetGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	id, err := store.CreateGroup(ctx, opt)
	require.NoError(t, err)

	// get non-exist feature group
	_, err = store.GetGroup(ctx, 0)
	require.Error(t, err)

	// get existing feature group by id
	group, err := store.GetGroup(ctx, id)
	require.NoError(t, err)
	require.Equal(t, opt.GroupName, group.Name)
	require.Equal(t, opt.EntityID, group.EntityID)
	require.Equal(t, opt.Description, group.Description)
	require.Equal(t, opt.Category, group.Category)

	// get existing feature group by name
	group, err = store.GetGroupByName(ctx, "device_info")
	require.NoError(t, err)
	require.Equal(t, opt.GroupName, group.Name)
	require.Equal(t, opt.EntityID, group.EntityID)
	require.Equal(t, opt.Description, group.Description)
	require.Equal(t, opt.Category, group.Category)
}

func TestListGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	deviceEntityID := prepareEntity(t, ctx, store, "device")
	userEntityID := prepareEntity(t, ctx, store, "user")

	deviceOpt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    deviceEntityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	userBaseOpt := metadata.CreateGroupOpt{
		GroupName:   "user_info",
		EntityID:    userEntityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	userBehaviorOpt := metadata.CreateGroupOpt{
		GroupName:   "user_profile",
		EntityID:    userEntityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	deviceGroupID, err := store.CreateGroup(ctx, deviceOpt)
	require.NoError(t, err)
	userGroupID, err := store.CreateGroup(ctx, userBaseOpt)
	require.NoError(t, err)
	_, err = store.CreateGroup(ctx, userBehaviorOpt)
	require.NoError(t, err)

	groups, err := store.ListGroup(ctx, nil, &[]int{})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(groups))

	groups, err = store.ListGroup(ctx, &deviceGroupID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(groups))

	groups, err = store.ListGroup(ctx, &userGroupID, nil)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(groups))

	groups, err = store.ListGroup(ctx, nil, nil)
	assert.NoError(t, err)
	assert.Equal(t, 3, len(groups))

	groups, err = store.ListGroup(ctx, intPtr(0), nil)
	assert.NoError(t, err)
	assert.Equal(t, 0, len(groups))

	ids := []int{deviceGroupID, userGroupID}
	groups, err = store.ListGroup(ctx, nil, &ids)
	assert.NoError(t, err)
	assert.Equal(t, 2, len(groups))

	groups, err = store.ListGroup(ctx, &userEntityID, &ids)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(groups))
}

func TestCreateGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}

	// create successfully
	groupID, err := store.CreateGroup(ctx, opt)
	require.NotEqual(t, 0, groupID)
	require.NoError(t, err)

	// cannot create feature group with same name
	_, err = store.CreateGroup(ctx, opt)
	require.Equal(t, fmt.Errorf("feature group device_info already exists"), err)

	// cannot create feature group with invalid category
	opt.Category = "invalid-category"
	_, err = store.CreateGroup(ctx, opt)
	require.NotNil(t, err)
}

func TestUpdateGroup(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	groupID, err := store.CreateGroup(ctx, opt)
	require.NoError(t, err)

	// update non-exist feature group
	require.NotNil(t, store.UpdateGroup(ctx, metadata.UpdateGroupOpt{
		GroupID: 0,
	}))

	// update existing feature group
	description := "new description"
	require.Nil(t, store.UpdateGroup(ctx, metadata.UpdateGroupOpt{
		GroupID:        groupID,
		NewDescription: &description,
	}))
}
