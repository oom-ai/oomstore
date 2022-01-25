package test_impl

import (
	"context"
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
			Description: "description",
		},
	})
	require.NoError(t, err)

	groupID, err := store.CreateGroup(ctx, metadata.CreateGroupOpt{
		GroupName:   "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.CategoryBatch,
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())
	return entityID, groupID
}

func TestCreateFeature(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		Description: "description",
		ValueType:   types.String,
	}

	_, err := store.CreateFeature(ctx, opt)
	assert.NoError(t, err)
}

func TestCreateFeatureWithSameName(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		ValueType:   types.String,
	}

	_, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)

	_, err = store.CreateFeature(ctx, opt)
	assert.Equal(t, "feature phone already exists", err.Error())
}

func TestCreateFeatureWithSQLKeyword(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "user",
		GroupID:     groupID,
		ValueType:   types.Int64,
		Description: "order",
	}

	_, err := store.CreateFeature(ctx, opt)
	assert.NoError(t, err)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	_, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "model",
		GroupID:     groupID,
	})
	assert.Error(t, err)
}

func TestGetFeature(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	id, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		Description: "description",
		ValueType:   types.String,
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
		ValueType:   types.String,
		Description: "description",
		GroupID:     1,
	}
	ignoreFeatureFields(feature)
	assert.Equal(t, expected, feature)
}

func TestGetFeatureByName(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	_, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		Description: "description",
		ValueType:   types.String,
	})
	require.NoError(t, err)

	// case 1: wrong feature name, return error
	_, err = store.GetFeatureByName(ctx, "g", "f")
	assert.EqualError(t, err, "feature group g not found")

	// case 2: correct feature name, return feature `phone`
	feature, err := store.GetFeatureByName(ctx, "device_info", "phone")
	assert.NoError(t, err)
	expected := &types.Feature{
		ID:          1,
		Name:        "phone",
		ValueType:   types.String,
		Description: "description",
		GroupID:     1,
	}
	ignoreFeatureFields(feature)
	assert.Equal(t, expected, feature)
}

func TestListCachedFeature(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	// case 1: no feature to list
	features := store.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{})
	assert.Equal(t, 0, features.Len())

	_, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		Description: "description",
		ValueType:   types.String,
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())

	// case 2: no condition, list all features
	features = store.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{})
	assert.Equal(t, 1, features.Len())

	// case 8: list features by FeatureNames
	features = store.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		FullNames: &[]string{"device_info.phone"},
	})
	assert.Equal(t, 1, len(features))
}

func TestListFeature(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

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
		Description: "description",
		ValueType:   types.String,
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

func TestUpdateFeature(t *testing.T, prepareStore PrepareStoreFn, destroyStore DestroyStoreFn) {
	t.Cleanup(destroyStore)

	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		FeatureName: "phone",
		GroupID:     groupID,
		Description: "description",
		ValueType:   types.String,
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
		ValueType:   types.String,
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
