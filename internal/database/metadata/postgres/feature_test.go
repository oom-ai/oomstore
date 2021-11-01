package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/stretchr/testify/assert"
)

func TestCreateFeature(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()
	ctx := context.Background()

	opt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device",
			DBValueType: "varchar(16)",
			Description: "description",
		},
		ValueType: "string",
	}

	err := db.CreateFeature(ctx, opt)
	assert.Nil(t, err)

	var feature types.Feature
	err = db.GetContext(ctx, &feature, "select * from feature where name = $1", opt.FeatureName)
	assert.Nil(t, err)
	assert.Equal(t, feature.Name, opt.FeatureName)
	assert.Equal(t, feature.GroupName, opt.GroupName)
	assert.Equal(t, feature.Description, opt.Description)
	assert.Equal(t, feature.DBValueType, opt.DBValueType)
	assert.Equal(t, feature.ValueType, opt.ValueType)
}

func TestCreateFeatureWithSameName(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()
	ctx := context.Background()

	opt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "phone",
			GroupName:   "device",
			DBValueType: "varchar(16)",
		},
	}

	err := db.CreateFeature(ctx, opt)
	assert.Nil(t, err)

	err = db.CreateFeature(ctx, opt)
	assert.Equal(t, err, fmt.Errorf("feature phone already exists"))

	opt.GroupName = "new group"
	assert.Equal(t, err, fmt.Errorf("feature phone already exists"))
}

func TestCreateFeatureWithSQLKeywrod(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()
	ctx := context.Background()

	opt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "group",
			GroupName:   "select",
			DBValueType: "int",
			Description: "order",
		},
	}

	err := db.CreateFeature(ctx, opt)
	assert.Nil(t, err)

	var feature types.Feature
	err = db.GetContext(ctx, &feature, "select * from feature where name = $1", "group")
	assert.Nil(t, err)
	assert.Equal(t, feature.Name, "group")
	assert.Equal(t, feature.GroupName, "select")
	assert.Equal(t, feature.Description, "order")
}

func TestCreateFeatureWithInvalidDataType(t *testing.T) {
	db := initAndOpenDB(t)
	defer db.Close()
	ctx := context.Background()

	opt := metadata.CreateFeatureOpt{
		CreateFeatureOpt: types.CreateFeatureOpt{
			FeatureName: "model",
			GroupName:   "phone",
			DBValueType: "invalid_type",
		},
	}

	err := db.CreateFeature(ctx, opt)
	assert.NotNil(t, err)
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

	rowsAffected, err := db.UpdateFeature(context.Background(), types.UpdateFeatureOpt{
		FeatureName: "invalid_feature_name",
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(0), rowsAffected)

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

	rowsAffected, err = db.UpdateFeature(context.Background(), types.UpdateFeatureOpt{
		FeatureName:    phoneOpt.FeatureName,
		NewDescription: "new description",
	})
	assert.Nil(t, err)
	assert.Equal(t, int64(1), rowsAffected)

	feature, err := db.GetFeature(context.Background(), "phone")
	assert.Nil(t, err)
	assert.Equal(t, "phone", feature.Name)
	assert.Equal(t, "device_baseinfo", feature.GroupName)
	assert.Equal(t, "varchar(16)", feature.DBValueType)
	assert.Equal(t, "new description", feature.Description)
}
