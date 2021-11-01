package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateEntity(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	assert.Equal(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}), fmt.Errorf("entity device already exists"))
}

func TestGetEntity(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	entity, err := db.GetEntity(context.Background(), "device")
	assert.Nil(t, err)
	assert.Equal(t, "device", entity.Name)
	assert.Equal(t, 32, entity.Length)
	assert.Equal(t, "description", entity.Description)

	entity, err = db.GetEntity(context.Background(), "invalid_entity_name")
	assert.Equal(t, err, sql.ErrNoRows)
	assert.Nil(t, entity)
}

func TestUpdateEntity(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	rowsAffected, err := db.UpdateEntity(context.Background(), types.UpdateEntityOpt{
		EntityName:     "device",
		NewDescription: "new description",
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	entity, err := db.GetEntity(context.Background(), "device")
	assert.Nil(t, err)
	assert.Equal(t, entity.Description, "new description")

	rowsAffected, err = db.UpdateEntity(context.Background(), types.UpdateEntityOpt{
		EntityName:     "invalid_entity_name",
		NewDescription: "new description",
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), rowsAffected)
}

func TestListEntity(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	entitys, err := db.ListEntity(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 0, len(entitys))

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))
	entitys, err = db.ListEntity(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 1, len(entitys))

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "user",
		Length:      16,
		Description: "description",
	}))
	entitys, err = db.ListEntity(context.Background())
	assert.Nil(t, err)
	assert.Equal(t, 2, len(entitys))
}
