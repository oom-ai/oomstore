package test_impl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareEntity(t *testing.T, ctx context.Context, store metadata.Store, name string) int {
	entityID, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  name,
			Description: "description",
		},
	})
	require.NoError(t, err)
	return entityID
}

func TestGetGroup(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.CategoryBatch,
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

func TestListGroup(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	deviceEntityID := prepareEntity(t, ctx, store, "device")
	userEntityID := prepareEntity(t, ctx, store, "user")

	deviceOpt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    deviceEntityID,
		Description: "description",
		Category:    types.CategoryBatch,
	}
	userBaseOpt := metadata.CreateGroupOpt{
		GroupName:   "user_info",
		EntityID:    userEntityID,
		Description: "description",
		Category:    types.CategoryBatch,
	}
	userBehaviorOpt := metadata.CreateGroupOpt{
		GroupName:   "user_profile",
		EntityID:    userEntityID,
		Description: "description",
		Category:    types.CategoryBatch,
	}
	deviceGroupID, err := store.CreateGroup(ctx, deviceOpt)
	require.NoError(t, err)
	userGroupID, err := store.CreateGroup(ctx, userBaseOpt)
	require.NoError(t, err)
	_, err = store.CreateGroup(ctx, userBehaviorOpt)
	require.NoError(t, err)

	groups, err := store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: nil,
		GroupIDs:  &[]int{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(groups))

	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: &[]int{deviceGroupID},
		GroupIDs:  nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(groups))

	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: &[]int{userGroupID},
		GroupIDs:  nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(groups))

	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: nil,
		GroupIDs:  nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(groups))

	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: &[]int{0},
		GroupIDs:  nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(groups))

	ids := []int{deviceGroupID, userGroupID}
	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: nil,
		GroupIDs:  &ids,
	})
	assert.NoError(t, err)
	assert.Equal(t, 2, len(groups))

	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: &[]int{userEntityID},
		GroupIDs:  &ids,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(groups))

	groups, err = store.ListGroup(ctx, metadata.ListGroupOpt{
		EntityIDs: &[]int{userEntityID, deviceEntityID},
		GroupIDs:  nil,
	})
	assert.NoError(t, err)
	assert.Equal(t, 3, len(groups))
}

func TestCreateGroup(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.CategoryBatch,
	}

	// create successfully
	groupID, err := store.CreateGroup(ctx, opt)
	require.NotEqual(t, 0, groupID)
	require.NoError(t, err)

	// cannot create feature group with same name
	_, err = store.CreateGroup(ctx, opt)
	require.Equal(t, "feature group device_info already exists", err.Error())

	// cannot create feature group with invalid category
	opt.Category = "invalid-category"
	_, err = store.CreateGroup(ctx, opt)
	require.NotNil(t, err)
}

func TestCreateStreamGroup(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "user")

	opt := metadata.CreateGroupOpt{
		GroupName:        "user-click",
		EntityID:         entityID,
		Description:      "description",
		Category:         types.CategoryStream,
		SnapshotInterval: 24 * int(time.Hour.Seconds()),
	}

	// create successfully
	groupID, err := store.CreateGroup(ctx, opt)
	require.NoError(t, err)
	require.NotEqual(t, 0, groupID)

	// cannot create feature group with same name
	_, err = store.CreateGroup(ctx, opt)
	require.Equal(t, "feature group user-click already exists", err.Error())

	// cannot create feature group with invalid category
	opt.Category = "invalid-category"
	_, err = store.CreateGroup(ctx, opt)
	require.NotNil(t, err)

	// cannot create stream group with snapshot_interval is zero
	opt.Category = types.CategoryStream
	opt.SnapshotInterval = 0
	_, err = store.CreateGroup(ctx, opt)
	require.Equal(t, fmt.Sprintf("the field SnapshotInterval of the stream group %s cannot be zero", opt.GroupName), err.Error())
}

func TestUpdateGroup(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()

	entityID := prepareEntity(t, ctx, store, "device")

	opt := metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.CategoryBatch,
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
