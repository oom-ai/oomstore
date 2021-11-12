package postgres_test

import (
	"context"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/metadatav2/postgres"
	"github.com/oom-ai/oomstore/internal/database/metadatav2/test"
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
	test.TestGetFeatureGroup(t, prepareStore)
}

func TestListFeatureGroup(t *testing.T) {
	test.TestListFeatureGroup(t, prepareStore)
}

func TestCreateFeatureGroup(t *testing.T) {
	test.TestCreateFeatureGroup(t, prepareStore)
}

func TestUpdateFeatureGroup(t *testing.T) {
	test.TestUpdateFeatureGroup(t, prepareStore)
}
