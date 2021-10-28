package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"sort"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/test"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
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

	if _, err := store.ExecContext(context.Background(), "drop database if exists oomstore"); err != nil {
		t.Fatal(err)
	}
	store.Close()

	if err := CreateDatabase(context.Background(), test.PostgresDbopt); err != nil {
		t.Fatal(err)
	}
}

func initAndOpenDB(t *testing.T) *DB {
	initDB(t)

	db, err := Open(&test.PostgresDbopt)
	if err != nil {
		t.Fatal(err)
	}
	return db
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
	}), fmt.Errorf("entity device already exist!"))
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

	assert.Nil(t, db.UpdateEntity(context.Background(), types.UpdateEntityOpt{
		EntityName:     "device",
		NewDescription: "new description",
	}))

	entity, err := db.GetEntity(context.Background(), "device")
	assert.Nil(t, err)
	assert.Equal(t, entity.Description, "new description")

	assert.Nil(t, db.UpdateEntity(context.Background(), types.UpdateEntityOpt{
		EntityName:     "invalid_entity_name",
		NewDescription: "new description",
	}))
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

func TestFeature(t *testing.T) {
	initDB(t)

	store, err := Open(&test.PostgresDbopt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

	assert.Nil(t, store.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	assert.Nil(t, store.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}))

	phoneOpt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}

	priceOpt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "price",
			GroupName:   "device",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}

	// test CreateFeature
	{
		errOpt := priceOpt
		errOpt.DBValueType = "varchar(16"
		assert.NotNil(t, store.CreateFeature(context.Background(), errOpt))

		assert.Nil(t, store.CreateFeature(context.Background(), phoneOpt))
		assert.Nil(t, store.CreateFeature(context.Background(), priceOpt))
	}

	// test GetFeature
	{
		feature, err := store.GetFeature(context.Background(), phoneOpt.FeatureName)
		assert.Nil(t, err)

		assert.Equal(t, phoneOpt.FeatureName, feature.Name)
		assert.Equal(t, phoneOpt.GroupName, feature.GroupName)
		assert.Equal(t, phoneOpt.DBValueType, feature.DBValueType)
		assert.Equal(t, phoneOpt.ValueType, feature.ValueType)
		assert.Equal(t, phoneOpt.Description, feature.Description)
		assert.Equal(t, "batch", feature.Category)
	}

	// testUpdateFeature
	{
		assert.Nil(t, store.UpdateFeature(context.Background(), types.UpdateFeatureOpt{
			FeatureName:    phoneOpt.FeatureName,
			NewDescription: "new description",
		}))

		feature, err := store.GetFeature(context.Background(), phoneOpt.FeatureName)
		assert.Nil(t, err)
		assert.Equal(t, "new description", feature.Description)
	}

	// test ListFeature
	{
		invalidGroupName := "invalid_group_name"
		features, err := store.ListFeature(context.Background(), types.ListFeatureOpt{GroupName: &invalidGroupName})
		assert.Nil(t, err)
		assert.Equal(t, 0, len(features))

		groupName := "device"
		features, err = store.ListFeature(context.Background(), types.ListFeatureOpt{GroupName: &groupName})
		assert.Nil(t, err)
		assert.Equal(t, 2, len(features))
	}
}

func TestFeatureGroup(t *testing.T) {
	initDB(t)

	store, err := Open(&test.PostgresDbopt)
	assert.Nil(t, err)
	defer store.Close()

	assert.Nil(t, store.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	baseInfoFg := metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "deviec_baseinfo",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}

	infoFg := metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "deviec_info",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}

	// test CreateFeatureGroup
	{
		errFg := baseInfoFg
		errFg.Category = "invalid-category"
		assert.NotNil(t, store.CreateFeatureGroup(context.Background(), errFg))
		assert.Nil(t, store.CreateFeatureGroup(context.Background(), baseInfoFg))
		assert.Nil(t, store.CreateFeatureGroup(context.Background(), infoFg))
	}

	// test GetFeatureGroup
	{
		fg, err := store.GetFeatureGroup(context.Background(), baseInfoFg.Name)
		assert.Nil(t, err)

		assert.Equal(t, baseInfoFg.Category, fg.Category)
		assert.Equal(t, baseInfoFg.EntityName, fg.EntityName)
		assert.Equal(t, baseInfoFg.Description, fg.Description)
		assert.Equal(t, baseInfoFg.Category, fg.Category)
	}

	// test UpdateFeatureGroup
	{
		description := "new description"
		revisionId := int32(2)
		assert.Nil(t, store.UpdateFeatureGroup(context.Background(), types.UpdateFeatureGroupOpt{
			GroupName:        baseInfoFg.Name,
			Description:      &description,
			OnlineRevisionId: &revisionId,
		}))

		fg, err := store.GetFeatureGroup(context.Background(), baseInfoFg.Name)
		assert.Nil(t, err)

		assert.Equal(t, "new description", fg.Description)
		assert.Equal(t, revisionId, *fg.OnlineRevisionID)
	}

	// test ListFeatureGroup
	{
		entityName := "device"
		fgs, err := store.ListFeatureGroup(context.Background(), &entityName)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(fgs))
	}
}

func TestRevision(t *testing.T) {
	initDB(t)

	store, err := Open(&test.PostgresDbopt)
	assert.Nil(t, err)
	defer store.Close()

	assert.Nil(t, store.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	assert.Nil(t, store.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "deviec_baseinfo",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}))

	opt1 := metadata.InsertRevisionOpt{
		GroupName:   "device_baseinfo",
		Revision:    20211028,
		DataTable:   "device_bastinfo_20211028",
		Description: "description",
	}

	opt2 := metadata.InsertRevisionOpt{
		GroupName:   "device_baseinfo",
		Revision:    20211029,
		DataTable:   "device_bastinfo_20211029",
		Description: "description",
	}

	// test InsertRevision
	{
		assert.Nil(t, store.InsertRevision(context.Background(), opt1))
		assert.Nil(t, store.InsertRevision(context.Background(), opt2))
	}

	// test GetRevision and GetRevisionsByDataTables
	{
		revision, err := store.GetRevision(context.Background(), metadata.GetRevisionOpt{
			GroupName: &opt1.GroupName,
			Revision:  &opt1.Revision,
		})
		assert.Nil(t, err)

		assert.Equal(t, opt1.GroupName, revision.GroupName)
		assert.Equal(t, opt1.Revision, revision.Revision)
		assert.Equal(t, opt1.DataTable, revision.DataTable)
		assert.Equal(t, opt1.Description, revision.Description)

		invalidGroupName := "invalid group name"
		invalidRevision := int64(0)
		revision, err = store.GetRevision(context.Background(), metadata.GetRevisionOpt{
			GroupName: &invalidGroupName,
			Revision:  &invalidRevision,
		})
		assert.NotNil(t, err)
		assert.Nil(t, revision)

		revisios, err := store.GetRevisionsByDataTables(context.Background(),
			[]string{"device_bastinfo_20211028", "device_bastinfo_20211029"})
		assert.Nil(t, err)

		assert.Equal(t, 2, len(revisios))
	}

	// test ListRevision
	{
		revisions, err := store.ListRevision(context.Background(), nil)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(revisions))

		groupName := "device_baseinfo"
		revisions, err = store.ListRevision(context.Background(), &groupName)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(revisions))
	}

}
