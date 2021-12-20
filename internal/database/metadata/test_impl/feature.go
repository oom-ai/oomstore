package test_impl

import (
	"context"
	"fmt"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareEntityAndGroup(t *testing.T, ctx context.Context, store metadata.Store) (int, int) {
	entityID, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: types.CreateEntityOpt{
			EntityName:  "device",
			Length:      32,
			Description: "description",
		},
	})
	require.NoError(t, err)

	groupID, err := store.CreateGroup(ctx, metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())
	return entityID, groupID
}

func TestCreateFeature(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}

	_, err := store.CreateFeature(ctx, opt)
	assert.NoError(t, err)
}

func TestCreateFeatureWithSameName(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
	}

	_, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)

	_, err = store.CreateFeature(ctx, opt)
	assert.Equal(t, err, fmt.Errorf("feature phone already exists"))
}

func TestCreateFeatureWithSQLKeyword(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "user",
		GroupID:     groupID,
		DBValueType: "int",
		Description: "order",
	}

	_, err := store.CreateFeature(ctx, opt)
	assert.NoError(t, err)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	_, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "model",
		GroupID:     groupID,
		DBValueType: "invalid_type",
	})
	assert.Error(t, err)
}

func TestGetFeature(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	id, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	// case 1: wrong featureID, return error
	_, err = store.GetFeature(ctx, 0)
	assert.EqualError(t, err, "feature 0 not found")

	// case 2: correct featureID, return feature `phone`
	feature, err := store.GetFeature(ctx, id)
	assert.NoError(t, err)
	expected := &types.Feature{
		ID:          1,
		Name:        "phone",
		ValueType:   types.STRING,
		DBValueType: "varchar(16)",
		Description: "description",
		GroupID:     1,
	}
	ignoreFeatureFields(feature)
	assert.Equal(t, expected, feature)
}

func TestGetFeatureByName(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	_, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	// case 1: wrong feature name, return error
	_, err = store.GetFeatureByName(ctx, "p")
	assert.EqualError(t, err, "feature p not found")

	// case 2: correct feature name, return feature `phone`
	feature, err := store.GetFeatureByName(ctx, "phone")
	assert.NoError(t, err)
	expected := &types.Feature{
		ID:          1,
		Name:        "phone",
		ValueType:   types.STRING,
		DBValueType: "varchar(16)",
		Description: "description",
		GroupID:     1,
	}
	ignoreFeatureFields(feature)
	assert.Equal(t, expected, feature)
}

func TestCacheListFeature(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	entityID, groupID := prepareEntityAndGroup(t, ctx, store)

	// case 1: no feature to list
	features := store.CacheListFeature(ctx, metadata.ListFeatureOpt{})
	assert.Equal(t, 0, features.Len())

	featureID, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())

	// case 2: no condition, list all features
	features = store.CacheListFeature(ctx, metadata.ListFeatureOpt{})
	assert.Equal(t, 1, features.Len())

	// case 3: list features by FeatureIDs
	features = store.CacheListFeature(ctx, metadata.ListFeatureOpt{
		FeatureIDs: &[]int{featureID},
	})
	assert.Equal(t, 1, features.Len())

	// case 4: list features by EntityID and FeatureIDs
	features = store.CacheListFeature(ctx, metadata.ListFeatureOpt{
		EntityID:   intPtr(entityID + 1),
		FeatureIDs: &[]int{featureID},
	})
	assert.Equal(t, 0, features.Len())

	// case 5: list features by GroupID and FeatureIDs
	features = store.CacheListFeature(ctx, metadata.ListFeatureOpt{
		GroupID:    intPtr(groupID + 1),
		FeatureIDs: &[]int{featureID},
	})
	assert.Equal(t, 0, features.Len())

	// case 6: list features by EntityID and empty FeatureIDs, return no feature
	features = store.CacheListFeature(ctx, metadata.ListFeatureOpt{
		EntityID:   &entityID,
		FeatureIDs: &[]int{},
	})
	assert.Equal(t, 0, len(features))

	// case 7: list features by EntityID
	features = store.CacheListFeature(ctx, metadata.ListFeatureOpt{
		EntityID: &entityID,
	})
	assert.Equal(t, 1, len(features))
}

func TestListFeature(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	entityID, groupID := prepareEntityAndGroup(t, ctx, store)

	// case 1: no feature to list
	features, err := store.ListFeature(ctx, metadata.ListFeatureOpt{})
	assert.NoError(t, err)
	assert.Equal(t, 0, features.Len())

	featureID, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	// case 2: no condition, list all features
	features, err = store.ListFeature(ctx, metadata.ListFeatureOpt{})
	assert.NoError(t, err)
	assert.Equal(t, 1, features.Len())

	// case 3: list features by FeatureIDs
	features, err = store.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureIDs: &[]int{featureID},
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, features.Len())

	// case 4: list features by EntityID and FeatureIDs
	features, err = store.ListFeature(ctx, metadata.ListFeatureOpt{
		EntityID:   intPtr(entityID + 1),
		FeatureIDs: &[]int{featureID},
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, features.Len())

	// case 5: list features by GroupID and FeatureIDs
	features, err = store.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID:    intPtr(groupID + 1),
		FeatureIDs: &[]int{featureID},
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, features.Len())

	// case 6: list features by EntityID and empty FeatureIDs, return no feature
	features, err = store.ListFeature(ctx, metadata.ListFeatureOpt{
		EntityID:   &entityID,
		FeatureIDs: &[]int{},
	})
	assert.NoError(t, err)
	assert.Equal(t, 0, len(features))

	// case 7: list features by EntityID
	features, err = store.ListFeature(ctx, metadata.ListFeatureOpt{
		EntityID: &entityID,
	})
	assert.NoError(t, err)
	assert.Equal(t, 1, len(features))
}

func TestUpdateFeature(t *testing.T, prepareStore PrepareStoreFn) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}
	id, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)

	// case 1: nothing to update
	err = store.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
		FeatureID: id + 1,
	})
	require.Error(t, err)

	// case 2: update description successfully
	err = store.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
		FeatureID:      id,
		NewDescription: stringPtr("new description"),
	})
	require.NoError(t, err)

	feature, err := store.GetFeature(ctx, id)
	assert.NoError(t, err)
	expected := &types.Feature{
		ID:          1,
		Name:        "phone",
		ValueType:   types.STRING,
		DBValueType: "varchar(16)",
		Description: "new description",
		GroupID:     1,
	}
	ignoreFeatureFields(feature)
	assert.Equal(t, expected, feature)
}

func ignoreFeatureFields(feature *types.Feature) {
	feature.CreateTime = time.Time{}
	feature.ModifyTime = time.Time{}
	feature.Group = nil
}
