package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareEntityAndGroup(t *testing.T, ctx context.Context, store metadatav2.Store) (int16, int16) {
	entityID, err := store.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	groupID, err := store.CreateFeatureGroup(ctx, metadatav2.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)
	require.NoError(t, store.Refresh())
	return entityID, groupID
}

func TestCreateFeature(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadatav2.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}

	_, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)
}

func TestCreateFeatureWithSameName(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadatav2.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
	}

	_, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)

	_, err = store.CreateFeature(ctx, opt)
	require.Equal(t, err, fmt.Errorf("feature phone already exists"))
}

func TestCreateFeatureWithSQLKeywrod(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadatav2.CreateFeatureOpt{
		Name:        "user",
		GroupID:     groupID,
		DBValueType: "int",
		Description: "order",
	}

	_, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	_, err := store.CreateFeature(ctx, metadatav2.CreateFeatureOpt{
		Name:        "model",
		GroupID:     groupID,
		DBValueType: "invalid_type",
	})
	require.Error(t, err)
}

func TestGetFeature(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	id, err := store.CreateFeature(ctx, metadatav2.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	_, err = store.GetFeature(ctx, 0)
	require.EqualError(t, err, "feature 0 not found")

	feature, err := store.GetFeature(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "phone", feature.Name)
	require.Equal(t, "device_info", feature.Group.Name)
	require.Equal(t, "varchar(16)", feature.DBValueType)
	require.Equal(t, "description", feature.Description)
}

func TestListFeature(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	entityID, groupID := prepareEntityAndGroup(t, ctx, store)

	features := store.ListFeature(ctx, metadatav2.ListFeatureOpt{})
	require.Equal(t, 0, features.Len())

	featureID, err := store.CreateFeature(ctx, metadatav2.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	features = store.ListFeature(ctx, metadatav2.ListFeatureOpt{})
	require.Equal(t, 1, features.Len())

	features = store.ListFeature(ctx, metadatav2.ListFeatureOpt{
		FeatureIDs: &[]int16{featureID},
	})
	require.Equal(t, 1, features.Len())

	features = store.ListFeature(ctx, metadatav2.ListFeatureOpt{
		EntityID:   int16Ptr(entityID + 1),
		FeatureIDs: &[]int16{featureID},
	})
	require.Equal(t, 0, features.Len())

	features = store.ListFeature(ctx, metadatav2.ListFeatureOpt{
		EntityID:   &entityID,
		FeatureIDs: &[]int16{},
	})
	require.Equal(t, 0, len(features))

	features = store.ListFeature(ctx, metadatav2.ListFeatureOpt{
		EntityID: &entityID,
	})
	require.Equal(t, 1, len(features))
}

func TestUpdateFeature(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadatav2.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}
	id, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)

	require.Error(t, store.UpdateFeature(ctx, metadatav2.UpdateFeatureOpt{
		FeatureID: id + 1,
	}))

	err = store.UpdateFeature(ctx, metadatav2.UpdateFeatureOpt{
		FeatureID:      id,
		NewDescription: "new description",
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	feature, err := store.GetFeature(ctx, id)
	require.NoError(t, err)
	require.Equal(t, "phone", feature.Name)
	require.Equal(t, "device_info", feature.Group.Name)
	require.Equal(t, "varchar(16)", feature.DBValueType)
	require.Equal(t, "new description", feature.Description)
}
