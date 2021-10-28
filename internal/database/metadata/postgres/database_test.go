package postgres

import (
	"context"
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
		features, err := store.ListFeature(context.Background(), &invalidGroupName)
		assert.Nil(t, err)
		assert.Equal(t, 0, len(features))

		groupName := "device"
		features, err = store.ListFeature(context.Background(), &groupName)
		assert.Nil(t, err)
		assert.Equal(t, 2, len(features))
	}
}

func TestRichFeature(t *testing.T) {
	initDB(t)

	store, err := Open(&test.PostgresDbopt)
	if err != nil {
		t.Fatal(err)
	}
	defer store.Close()

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
	assert.Nil(t, store.CreateFeature(context.Background(), phoneOpt))
	assert.Nil(t, store.CreateFeature(context.Background(), priceOpt))

	// test GetRichFeature
	{
		feature, err := store.GetRichFeature(context.Background(), "phone")
		assert.Nil(t, err)

		assert.Equal(t, phoneOpt.FeatureName, feature.Name)
		assert.Equal(t, phoneOpt.GroupName, feature.GroupName)
		assert.Equal(t, phoneOpt.DBValueType, feature.DBValueType)
		assert.Equal(t, phoneOpt.ValueType, feature.ValueType)
		assert.Equal(t, phoneOpt.Description, feature.Description)
		assert.Equal(t, "batch", feature.Category)

		assert.Nil(t, feature.Revision)
		assert.Nil(t, feature.DataTable)
	}

	// test GetRichFeatures
	{
		features, err := store.GetRichFeatures(context.Background(), []string{phoneOpt.FeatureName, priceOpt.FeatureName})
		assert.Nil(t, err)

		assert.Equal(t, 2, len(features))
	}

	// test ListRichFeatuer
	{
		groupName := "device"
		features, err := store.ListRichFeature(context.Background(), types.ListFeatureOpt{GroupName: &groupName})
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
		revision := int64(20211028)
		assert.Nil(t, store.UpdateFeatureGroup(context.Background(), types.UpdateFeatureGroupOpt{
			GroupName:   baseInfoFg.Name,
			Description: &description,
			Revision:    &revision,
		}))

		fg, err := store.GetFeatureGroup(context.Background(), baseInfoFg.Name)
		assert.Nil(t, err)

		assert.Equal(t, "new description", fg.Description)
		assert.Equal(t, revision, *fg.Revision)
	}

	// test ListFeatureGroup
	{
		entityName := "device"
		fgs, err := store.ListFeatureGroup(context.Background(), &entityName)
		assert.Nil(t, err)

		assert.Equal(t, 2, len(fgs))
	}
}
