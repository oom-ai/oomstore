package postgres_test

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/metadata/postgres"
	"github.com/oom-ai/oomstore/internal/database/metadata/test"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/require"
)

func prepareEntityAndGroup(t *testing.T, ctx context.Context, db *postgres.DB) (int16, int16) {
	entityID, err := db.CreateEntity(ctx, metadata.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	groupID, err := db.CreateFeatureGroup(ctx, metadata.CreateFeatureGroupOpt{
		Name:        "device_info",
		EntityID:    entityID,
		Description: "description",
		Category:    types.BatchFeatureCategory,
	})
	require.NoError(t, err)
	require.NoError(t, db.Refresh())
	return entityID, groupID
}

func TestCreateFeature(t *testing.T) {
	test.TestCreateFeature(t, prepareStore)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	test.TestCreateFeatureWithSameName(t, prepareStore)
}

func TestCreateFeatureWithSQLKeywrod(t *testing.T) {
	test.TestCreateFeatureWithSQLKeywrod(t, prepareStore)
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	test.TestCreateFeatureWithInvalidDataType(t, prepareStore)
}

func TestGetFeature(t *testing.T) {
	test.TestGetFeature(t, prepareStore)
}

func TestListFeature(t *testing.T) {
	test.TestListFeature(t, prepareStore)
}

func TestUpdateFeature(t *testing.T) {
	test.TestUpdateFeature(t, prepareStore)
}
