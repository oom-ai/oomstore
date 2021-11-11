package postgres_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/metadatav2/postgres"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

// create an entity with given name
func prepareEntity(t *testing.T, ctx context.Context, db *postgres.DB, name string) int16 {
	entityId, err := db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        name,
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)
	return entityId
}

func TestGetFeatureGroup(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	entityId := prepareEntity(t, ctx, db, "device")

	opt := metadatav2.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	id, err := db.CreateFeatureGroup(ctx, opt)
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	// get non-exist feature group
	_, err = db.GetFeatureGroup(ctx, 0)
	assert.Error(t, err)

	// get existing feature group
	featureGroup, err := db.GetFeatureGroup(ctx, id)
	assert.NoError(t, err)
	assert.Equal(t, opt.Name, featureGroup.Name)
	assert.Equal(t, opt.EntityID, featureGroup.EntityID)
	assert.Equal(t, opt.Description, featureGroup.Description)
	assert.Equal(t, opt.Category, featureGroup.Category)
}

func TestListFeatureGroup(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	deviceEntityId := prepareEntity(t, ctx, db, "device")
	userEntityId := prepareEntity(t, ctx, db, "user")

	deviceOpt := metadatav2.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    deviceEntityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	userBaseOpt := metadatav2.CreateFeatureGroupOpt{
		Name:        "user_baseinfo",
		EntityID:    userEntityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	userBehaviorOpt := metadatav2.CreateFeatureGroupOpt{
		Name:        "user_behaviorinfo",
		EntityID:    userEntityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	deviceGroupID, err := db.CreateFeatureGroup(ctx, deviceOpt)
	require.NoError(t, err)
	userGroupID, err := db.CreateFeatureGroup(ctx, userBaseOpt)
	require.NoError(t, err)
	_, err = db.CreateFeatureGroup(ctx, userBehaviorOpt)
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	assert.Equal(t, 1, len(db.ListFeatureGroup(ctx, &deviceGroupID)))
	assert.Equal(t, 2, len(db.ListFeatureGroup(ctx, &userGroupID)))
	assert.Equal(t, 3, len(db.ListFeatureGroup(ctx, nil)))
	assert.Equal(t, 0, len(db.ListFeatureGroup(ctx, int16Ptr(0))))
}

func TestCreateFeatureGroup(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	entityId := prepareEntity(t, ctx, db, "device")

	opt := metadatav2.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}

	// create successfully
	featureGroupId, err := db.CreateFeatureGroup(ctx, opt)
	assert.NotEqual(t, int16(0), featureGroupId)
	assert.NoError(t, err)

	// cannot create feature group with same name
	_, err = db.CreateFeatureGroup(ctx, opt)
	assert.Equal(t, fmt.Errorf("feature group device_baseinfo already exists"), err)

	// cannot create feature group with invalid category
	opt.Category = "invalid-category"
	_, err = db.CreateFeatureGroup(ctx, opt)
	assert.NotNil(t, err)
}

func TestUpdateFeatureGroup(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	entityId := prepareEntity(t, ctx, db, "device")

	opt := metadatav2.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	featureGroupId, err := db.CreateFeatureGroup(ctx, opt)
	require.NoError(t, err)

	// update non-exist feature group
	assert.NotNil(t, db.UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
		GroupID: 0,
	}))

	// update existing feature group
	description := "new description"
	assert.Nil(t, db.UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
		GroupID:     featureGroupId,
		Description: &description,
	}))
}
