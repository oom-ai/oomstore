package _init

import (
	"context"
	"log"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	DBOption database.Option
}

func Init(ctx context.Context, option *Option) {
	if err := database.CreateDatabase(ctx, &option.DBOption); err != nil {
		log.Fatalf("failed creating database: %v", err)
	}

	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	if err := database.CreateFeatureConfigTable(ctx, db); err != nil {
		log.Fatalf("failed initializing feature store: %v", err)
	}

	if err := database.CreateFeatureRevisionTable(ctx, db); err != nil {
		log.Fatalf("failed initializing feature store: %v", err)
	}

	log.Println("succeeded.")
}
