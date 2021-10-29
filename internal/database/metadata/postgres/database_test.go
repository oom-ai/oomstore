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

func TestCreateFeature(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	phoneOpt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}

	assert.Nil(t, db.CreateFeature(context.Background(), phoneOpt))
	assert.Equal(t, db.CreateFeature(context.Background(), phoneOpt), fmt.Errorf("feature phone already exist"))
}

func TestGetFeature(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	feature, err := db.GetFeature(context.Background(), "invalid_feature_name")
	assert.NotNil(t, err)
	assert.Nil(t, feature)

	assert.Nil(t, db.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}))
	assert.Nil(t, db.CreateFeature(context.Background(), metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}))

	feature, err = db.GetFeature(context.Background(), "phone")
	assert.Nil(t, err)
	assert.Equal(t, "phone", feature.Name)
	assert.Equal(t, "device", feature.GroupName)
	assert.Equal(t, "varchar(16)", feature.DBValueType)
	assert.Equal(t, "description", feature.Description)
}

func TestListFeature(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	features, err := db.ListFeature(context.Background(), types.ListFeatureOpt{})
	assert.Nil(t, err)
	assert.Equal(t, 0, features.Len())

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))

	assert.Nil(t, db.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device_baseinfo",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}))

	assert.Nil(t, db.CreateFeature(context.Background(), metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device_baseinfo",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}))

	features, err = db.ListFeature(context.Background(), types.ListFeatureOpt{})
	assert.Nil(t, err)
	assert.Equal(t, 1, features.Len())

	features, err = db.ListFeature(context.Background(), types.ListFeatureOpt{
		FeatureNames: []string{"phone", "model"},
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, features.Len())

	entityName := "invalid_entity_name"
	features, err = db.ListFeature(context.Background(), types.ListFeatureOpt{
		EntityName:   &entityName,
		FeatureNames: []string{"phone", "model"},
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, features.Len())

	entityName = "device"
	features, err = db.ListFeature(context.Background(), types.ListFeatureOpt{
		EntityName:   &entityName,
		FeatureNames: []string{},
	})
	assert.Nil(t, err)
	assert.Equal(t, 0, len(features))

	entityName = "device"
	features, err = db.ListFeature(context.Background(), types.ListFeatureOpt{
		EntityName: &entityName,
	})
	assert.Nil(t, err)
	assert.Equal(t, 1, len(features))
}

func TestUpdateFeatuer(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	assert.Nil(t, db.UpdateFeature(context.Background(), types.UpdateFeatureOpt{
		FeatureName: "invalid_feature_name",
	}))

	assert.Nil(t, db.CreateEntity(context.Background(), types.CreateEntityOpt{
		Name:        "device",
		Length:      32,
		Description: "description",
	}))
	assert.Nil(t, db.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device_baseinfo",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}))

	phoneOpt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device_baseinfo",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}
	assert.Nil(t, db.CreateFeature(context.Background(), phoneOpt))

	assert.Nil(t, db.UpdateFeature(context.Background(), types.UpdateFeatureOpt{
		FeatureName:    phoneOpt.FeatureName,
		NewDescription: "new description",
	}))

	feature, err := db.GetFeature(context.Background(), "phone")
	assert.Nil(t, err)
	assert.Equal(t, "phone", feature.Name)
	assert.Equal(t, "device_baseinfo", feature.GroupName)
	assert.Equal(t, "varchar(16)", feature.DBValueType)
	assert.Equal(t, "new description", feature.Description)
}

func TestCreateFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	infoFg := metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device_info",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}

	assert.Nil(t, db.CreateFeatureGroup(context.Background(), infoFg))
	assert.Equal(t, db.CreateFeatureGroup(context.Background(), infoFg), fmt.Errorf("feature group device_info already exist"))

	errInfoFg := infoFg
	errInfoFg.Category = "invalid-category"
	assert.NotNil(t, db.CreateFeatureGroup(context.Background(), errInfoFg))
}

func TestGetFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	fg, err := db.GetFeatureGroup(context.Background(), "invalid-feature-group")
	assert.NotNil(t, err)
	assert.Nil(t, fg)

	infoFg := metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device_info",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}
	assert.Nil(t, db.CreateFeatureGroup(context.Background(), infoFg))

	fg, err = db.GetFeatureGroup(context.Background(), "device_info")
	assert.Nil(t, err)
	assert.Equal(t, infoFg.Category, fg.Category)
	assert.Equal(t, infoFg.EntityName, fg.EntityName)
	assert.Equal(t, infoFg.Description, fg.Description)
	assert.Equal(t, infoFg.Category, fg.Category)

	fg, err = db.GetFeatureGroup(context.Background(), "invalid-feature-group")
	assert.NotNil(t, err)
	assert.Nil(t, fg)
}

func TestUpdateFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	rowsAffected, err := db.UpdateFeatureGroup(context.Background(), types.UpdateFeatureGroupOpt{
		GroupName: "invalid-group",
	})
	assert.NotNil(t, err)
	assert.Equal(t, int64(0), rowsAffected)

	description := "new description"
	rowsAffected, err = db.UpdateFeatureGroup(context.Background(), types.UpdateFeatureGroupOpt{
		GroupName:   "invalid-group",
		Description: &description,
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), rowsAffected)
}

func TestListFeatureGroup(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()

	entityName := "invalid-entity-name"
	fgs, err := db.ListFeatureGroup(context.Background(), &entityName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(fgs))

	assert.Nil(t, db.CreateFeatureGroup(context.Background(), metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: types.CreateFeatureGroupOpt{
			Name:        "device_baseinfo",
			EntityName:  "device",
			Description: "description",
		},
		Category: types.BatchFeatureCategory,
	}))

	entityName = "device"
	fgs, err = db.ListFeatureGroup(context.Background(), &entityName)
	assert.Nil(t, err)
	assert.Equal(t, 1, len(fgs))

	entityName = "invalid_entity_name"
	fgs, err = db.ListFeatureGroup(context.Background(), &entityName)
	assert.Nil(t, err)
	assert.Equal(t, 0, len(fgs))
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
