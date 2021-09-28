package describe_feature

import (
	"context"
	"fmt"
	"log"
	"os"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Name     string
	Group    string
	DBOption database.Option
}

func Run(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	record, err := db.GetFeatureConfig(ctx, option.Group, option.Name)
	if err != nil {
		log.Fatalf("failed querying feature config, group_name=%s, feature_name=%s: %v", option.Group, option.Name, err)
	}
	if record == nil {
		fmt.Fprintln(os.Stderr, "Feature not found.")
		os.Exit(4)
	}

	fmt.Println(record.String())
}
