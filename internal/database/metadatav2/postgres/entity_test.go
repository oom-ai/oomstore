package postgres

import (
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCreateEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	opt := types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}
	_, err := db.CreateEntity(ctx, opt)
	assert.NoError(t, err)

	_, err = db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	assert.Equal(t, err, fmt.Errorf("entity device already exists"))
}

func TestGetEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	opt := types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}

	_, err := db.CreateEntity(ctx, opt)
	require.NoError(t, err)

	require.NoError(t, db.Refresh())

	entity := db.GetEntity(ctx, opt.Name)
	require.NotNil(t, entity)
	assert.Equal(t, opt.Name, entity.Name)
	assert.Equal(t, opt.Length, entity.Length)
	assert.Equal(t, opt.Description, entity.Description)

	assert.Nil(t, db.GetEntity(ctx, "invalid_entity_name"))
}

func TestUpdateEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	_, err := db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	require.NoError(t, err)

	rowsAffected, err := db.UpdateEntity(ctx, types.UpdateEntityOpt{
		EntityName:     "device",
		NewDescription: "new description",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	require.NoError(t, db.Refresh())

	entity := db.GetEntity(ctx, "device")
	require.NotNil(t, entity)
	assert.Equal(t, entity.Description, "new description")

	rowsAffected, err = db.UpdateEntity(ctx, types.UpdateEntityOpt{
		EntityName:     "invalid_entity_name",
		NewDescription: "new description",
	})
	require.NoError(t, err)
	assert.Equal(t, int64(0), rowsAffected)
}

func TestListEntity(t *testing.T) {
	ctx, db := prepareStore(t)
	defer db.Close()

	require.NoError(t, db.Refresh())

	entitys := db.ListEntity(ctx)
	assert.Equal(t, 0, len(entitys))

	_, err := db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	})
	assert.NoError(t, err)

	require.NoError(t, db.Refresh())

	entitys = db.ListEntity(ctx)
	assert.Equal(t, 1, len(entitys))
	_, err = db.CreateEntity(ctx, types.CreateEntityOpt{
		Name:        "user",
		Length:      16,
		Description: "description",
	})
	assert.NoError(t, err)

	require.NoError(t, db.Refresh())

	entitys = db.ListEntity(ctx)
	assert.Equal(t, 2, len(entitys))
}
