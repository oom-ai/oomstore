package postgres_test

import (
	"fmt"
	"testing"

	metadatav2 "github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	opt := metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}
	_, err := db.CreateEntity(ctx, opt)
	assert.NoError(t, err)

	_, err = db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	assert.Equal(t, err, fmt.Errorf("entity device already exists"))
}

func TestGetEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	opt := metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}

	id, err := db.CreateEntity(ctx, opt)
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	entity, err := db.GetEntity(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, opt.Name, entity.Name)
	assert.Equal(t, opt.Length, entity.Length)
	assert.Equal(t, opt.Description, entity.Description)

	_, err = db.GetEntity(ctx, 0)
	assert.EqualError(t, err, "feature entity 0 not found")
}

func TestUpdateEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	id, err := db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	require.NoError(t, db.UpdateEntity(ctx, metadatav2.UpdateEntityOpt{
		EntityID:       id,
		NewDescription: "new description",
	}))

	require.NoError(t, db.Refresh())

	entity, err := db.GetEntity(ctx, id)
	require.NoError(t, err)
	assert.Equal(t, entity.Description, "new description")

	require.Error(t, db.UpdateEntity(ctx, metadatav2.UpdateEntityOpt{
		EntityID:       id + 1,
		NewDescription: "new description",
	}))
}

func TestListEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	require.NoError(t, db.Refresh())

	entitys := db.ListEntity(ctx)
	assert.Equal(t, 0, len(entitys))

	_, err := db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	assert.NoError(t, err)

	require.NoError(t, db.Refresh())

	entitys = db.ListEntity(ctx)
	assert.Equal(t, 1, len(entitys))
	_, err = db.CreateEntity(ctx, metadatav2.CreateEntityOpt{
		Name:        "user",
		Length:      16,
		Description: "description",
	})
	assert.NoError(t, err)

	require.NoError(t, db.Refresh())

	entitys = db.ListEntity(ctx)
	assert.Equal(t, 2, len(entitys))
}
