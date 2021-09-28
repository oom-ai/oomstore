package register_feature

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	database.FeatureConfig
	DBOption database.Option
}

func Run(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	log.Println("obtainning value type...")
	valueType, err := database.GetFeatureValueType(ctx, db, &option.FeatureConfig)
	if err != nil {
		log.Fatalf("failed obtainning value type: %v", err)
	}
	option.ValueType = valueType

	log.Println("registering new feature...")
	if err = database.RegisterFeatureConfig(ctx, db, option.FeatureConfig); err != nil {
		log.Fatalf("failed registering a new feature: %v", err)
	}

	log.Println("succeeded.")
}
