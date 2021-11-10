package postgres_test

import (
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFeatureGroup(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	entityId, err := db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	opt := metadatav2.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}

	// create successfully
	var featureGroupId int16
	featureGroupId, err = db.CreateFeatureGroup(ctx, opt)
	assert.NotEqual(t, int16(0), featureGroupId)
	assert.NoError(t, err)

	// cannot create feature group with same name
	featureGroupId, err = db.CreateFeatureGroup(ctx, opt)
	assert.Equal(t, int16(0), featureGroupId)
	assert.Equal(t, fmt.Errorf("feature group device_baseinfo already exists"), err)

	// cannot create feature group with invalid category
	opt.Category = "invalid-category"
	featureGroupId, err = db.CreateFeatureGroup(ctx, opt)
	assert.Equal(t, int16(0), featureGroupId)
	assert.NotNil(t, err)
}

func TestUpdateFeatureGroup(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	entityId, err := db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	opt := metadatav2.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityID:    entityId,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	var featureGroupId int16
	featureGroupId, err = db.CreateFeatureGroup(ctx, opt)
	require.NoError(t, err)

	// update non-exist feature group
	assert.NotNil(t, db.UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
		GroupID: int16(0),
	}))

	// update existing feature group
	description := "new description"
	assert.Nil(t, db.UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
		GroupID:     int16(featureGroupId),
		Description: &description,
	}))
}
