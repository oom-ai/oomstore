package postgres

import (
	"context"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/onestore-ai/onestore/internal/database/test"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func initDB(t *testing.T) {
	opt := test.PostgresDbopt
	store, err := Open(&types.PostgresDbOpt{
		Host:     opt.Host,
		Port:     opt.Port,
		User:     opt.User,
		Pass:     opt.Pass,
		Database: "test",
	})
	if err != nil {
		t.Fatal(err)
	}

	if _, err := store.ExecContext(context.Background(), "drop database if exists onestore"); err != nil {
		t.Fatal(err)
	}
	store.Close()

	if err := CreateDatabase(context.Background(), test.PostgresDbopt); err != nil {
		t.Fatal(err)
	}
}

func TestCreateDatabase(t *testing.T) {
	ctx := context.Background()
	if err := CreateDatabase(ctx, test.PostgresDbopt); err != nil {
		t.Fatal(err)
	}

	store, err := Open(&test.PostgresDbopt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	var tables []string
	if err = store.SelectContext(ctx, &tables,
		`SELECT table_name
FROM information_schema.tables
WHERE table_schema = 'public'
ORDER BY table_name;`); err != nil {
		t.Fatal(err)
	}

	var wantTables []string
	for table := range META_TABLE_SCHEMAS {
		wantTables = append(wantTables, table)
	}
	for table := range META_VIEW_SCHEMAS {
		wantTables = append(wantTables, table)
	}

	sort.Slice(tables, func(i, j int) bool {
		return tables[i] < tables[j]
	})
	sort.Slice(wantTables, func(i, j int) bool {
		return wantTables[i] < wantTables[j]
	})
	assert.Equal(t, wantTables, tables)
}

func TestEntity(t *testing.T) {
	initDB(t)

	store, err := Open(&test.PostgresDbopt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	deviceEntity := types.Entity{
		Name:        "device",
		Length:      32,
		Description: "description",
	}
	userEntity := types.Entity{
		Name:        "user",
		Length:      64,
		Description: "description",
	}

	// test CreateEntity
	{
		if err := store.CreateEntity(context.Background(), types.CreateEntityOpt{
			Name:        deviceEntity.Name,
			Length:      deviceEntity.Length,
			Description: deviceEntity.Description,
		}); err != nil {
			t.Error(err)
		}

		if err := store.CreateEntity(context.Background(), types.CreateEntityOpt{
			Name:        userEntity.Name,
			Length:      userEntity.Length,
			Description: userEntity.Description,
		}); err != nil {
			t.Error(err)
		}
	}

	// test GetEntity
	{
		entity, err := store.GetEntity(context.Background(), "device")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, deviceEntity.Name, entity.Name)
		assert.Equal(t, deviceEntity.Length, entity.Length)
		assert.Equal(t, deviceEntity.Description, entity.Description)
	}

	// test UpdateEntity
	{
		if err := store.UpdateEntity(context.Background(), types.UpdateEntityOpt{
			EntityName:     "user",
			NewDescription: "new description",
		}); err != nil {
			t.Error(err)
		}

		entity, err := store.GetEntity(context.Background(), "user")
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, "new description", entity.Description)
	}

	// test ListEntity
	{
		entitys, err := store.ListEntity(context.Background())
		if err != nil {
			t.Error(err)
		}
		assert.Equal(t, 2, len(entitys))
	}
}
