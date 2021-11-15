package test

import (
	"context"
	"fmt"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func prepareEntityAndGroup(t *testing.T, ctx context.Context, store metadata.Store) (int16, int16) {
	entityID, err := store.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	groupID, err := store.CreateFeatureGroup(ctx, metadata.CreateFeatureGroupOpt{
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

	opt := metadata.CreateFeatureOpt{
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

	opt := metadata.CreateFeatureOpt{
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

	opt := metadata.CreateFeatureOpt{
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

	_, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
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

	id, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
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

	features := store.ListFeature(ctx, metadata.ListFeatureOpt{})
	require.Equal(t, 0, features.Len())

	featureID, err := store.CreateFeature(ctx, metadata.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	})
	require.NoError(t, err)

	require.NoError(t, store.Refresh())

	features = store.ListFeature(ctx, metadata.ListFeatureOpt{})
	require.Equal(t, 1, features.Len())

	features = store.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureIDs: &[]int16{featureID},
	})
	require.Equal(t, 1, features.Len())

	features = store.ListFeature(ctx, metadata.ListFeatureOpt{
		EntityID:   int16Ptr(entityID + 1),
		FeatureIDs: &[]int16{featureID},
	})
	require.Equal(t, 0, features.Len())

	features = store.ListFeature(ctx, metadata.ListFeatureOpt{
		EntityID:   &entityID,
		FeatureIDs: &[]int16{},
	})
	require.Equal(t, 0, len(features))

	features = store.ListFeature(ctx, metadata.ListFeatureOpt{
		EntityID: &entityID,
	})
	require.Equal(t, 1, len(features))
}

func TestUpdateFeature(t *testing.T, prepareStore PrepareStoreRuntimeFunc) {
	ctx, store := prepareStore(t)
	defer store.Close()
	_, groupID := prepareEntityAndGroup(t, ctx, store)

	opt := metadata.CreateFeatureOpt{
		Name:        "phone",
		GroupID:     groupID,
		DBValueType: "varchar(16)",
		Description: "description",
		ValueType:   "string",
	}
	id, err := store.CreateFeature(ctx, opt)
	require.NoError(t, err)

	require.Error(t, store.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
		FeatureID: id + 1,
	}))

	err = store.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
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
