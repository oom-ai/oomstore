package postgres_test

import (
	"context"
	"fmt"
	"testing"

	metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/metadatav2/postgres"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func prepareEntityAndGroup(t *testing.T, ctx context.Context, db *postgres.DB) (int16, int16) {
	entityID, err := db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	groupID, err := db.CreateFeatureGroup(ctx, metadatav2.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)

	return entityID, groupID
}

func TestCreateFeature(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	opt := metadatav2.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}

	_, err := db.CreateFeature(ctx, opt)
	require.NoError(t, err)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	opt := metadatav2.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
	}

	_, err := db.CreateFeature(ctx, opt)
	require.NoError(t, err)

	_, err = db.CreateFeature(ctx, opt)
	assert.Equal(t, err, fmt.Errorf("feature phone already exists"))
}

func TestCreateFeatureWithSQLKeywrod(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	opt := metadatav2.CreateFeatureOpt{
		FeatureName: "user",
		GroupID:     groupID,
		DBValueType: "int",
		Description: "order",
	}

	_, err := db.CreateFeature(ctx, opt)
	require.NoError(t, err)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	_, err := db.CreateFeature(ctx, metadatav2.CreateFeatureOpt{
		FeatureName: "model",
		GroupID:     groupID,
		DBValueType: "invalid_type",
	})
	require.Error(t, err)
}

func TestGetFeature(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	feature := db.GetFeature(ctx, "invalid_feature_name")
	assert.Nil(t, feature)

	_, err := db.CreateFeature(ctx, metadatav2.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	feature = db.GetFeature(ctx, "phone")
	require.NotNil(t, feature)
	assert.Equal(t, "phone", feature.Name)
	assert.Equal(t, "device_info", feature.Group.Name)
	assert.Equal(t, "varchar(16)", feature.DBValueType)
	assert.Equal(t, "description", feature.Description)
}

func TestListFeature(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	features := db.ListFeature(ctx, metadatav2.ListFeatureOpt{})
	assert.Equal(t, 0, features.Len())

	_, err := db.CreateFeature(ctx, metadatav2.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	features = db.ListFeature(ctx, metadatav2.ListFeatureOpt{})
	assert.Equal(t, 1, features.Len())

	features = db.ListFeature(ctx, metadatav2.ListFeatureOpt{
		FeatureNames: []string{"phone", "model"},
	})
	assert.Equal(t, 1, features.Len())

	entityName := "invalid_entity_name"
	features = db.ListFeature(ctx, metadatav2.ListFeatureOpt{
		EntityName:   &entityName,
		FeatureNames: []string{"phone", "model"},
	})
	assert.Equal(t, 0, features.Len())

	entityName = "device"
	features = db.ListFeature(ctx, metadatav2.ListFeatureOpt{
		EntityName:   &entityName,
		FeatureNames: []string{},
	})
	assert.Equal(t, 0, len(features))

	entityName = "device"
	features = db.ListFeature(ctx, metadatav2.ListFeatureOpt{
		EntityName: &entityName,
	})
	assert.Equal(t, 1, len(features))
}

func TestUpdateFeature(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, db)

	opt := metadatav2.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}
	id, err := db.CreateFeature(ctx, opt)
	require.NoError(t, err)

	require.Error(t, db.UpdateFeature(ctx, metadatav2.UpdateFeatureOpt{
		FeatureID: id + 1,
	}))

	err = db.UpdateFeature(ctx, metadatav2.UpdateFeatureOpt{
		FeatureID:      id,
		NewDescription: "new description",
	})
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	feature := db.GetFeature(ctx, "phone")
	require.NotNil(t, feature)
	assert.Equal(t, "phone", feature.Name)
	assert.Equal(t, "device_baseinfo", feature.Group.Name)
	assert.Equal(t, "varchar(16)", feature.DBValueType)
	assert.Equal(t, "new description", feature.Description)
}
