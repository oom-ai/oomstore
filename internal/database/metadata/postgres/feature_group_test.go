package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	infoFg := metadata.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityId:    1,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}

	assert.Nil(t, db.CreateFeatureGroup(context.Background(), infoFg))
	assert.Equal(t, db.CreateFeatureGroup(context.Background(), infoFg), fmt.Errorf("feature group device_info already exist"))

	errInfoFg := infoFg
	errInfoFg.Category = "invalid-category"
	assert.NotNil(t, db.CreateFeatureGroup(context.Background(), errInfoFg))
}

func TestGetFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	fg, err := db.GetFeatureGroup(context.Background(), "invalid-feature-group")
	assert.NotNil(t, err)
	assert.Nil(t, fg)

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))
	entity, err := db.GetEntity(context.Background(), "device")
	require.NoError(t, err)

	infoFg := metadata.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityId:    entity.ID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}
	assert.Nil(t, db.CreateFeatureGroup(context.Background(), infoFg))

	fg, err = db.GetFeatureGroup(context.Background(), "device_info")
	assert.Nil(t, err)
	assert.Equal(t, infoFg.Category, fg.Category)
	assert.Equal(t, infoFg.EntityId, fg.EntityId)
	assert.Equal(t, infoFg.Description, fg.Description)
	assert.Equal(t, infoFg.Category, fg.Category)

	fg, err = db.GetFeatureGroup(context.Background(), "invalid-feature-group")
	assert.NotNil(t, err)
	assert.Nil(t, fg)
}

func TestUpdateFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	rowsAffected, err := db.UpdateFeatureGroup(context.Background(), types.UpdateFeatureGroupOpt{
		GroupName: "invalid-group",
	})
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), rowsAffected)

	description := "new description"
	rowsAffected, err = db.UpdateFeatureGroup(context.Background(), types.UpdateFeatureGroupOpt{
		GroupName:   "invalid-group",
		Description: &description,
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), rowsAffected)
}

func TestListFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	entityName := "invalid-entity-name"
	fgs, err := db.ListFeatureGroup(context.Background(), &entityName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(fgs))

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))
	entity, err := db.GetEntity(context.Background(), "device")
	require.NoError(t, err)

	assert.NoError(t, db.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		Name:        "device_baseinfo",
		EntityId:    entity.ID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	}))

	fgs, err = db.ListFeatureGroup(context.Background(), &entity.Name)
	assert.NoError(t, err)
	assert.Equal(t, 1, len(fgs))

	entityName = "invalid_entity_name"
	fgs, err = db.ListFeatureGroup(context.Background(), &entityName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(fgs))
}
