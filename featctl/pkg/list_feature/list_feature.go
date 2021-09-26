package list_feature

import (
	"context"
	"fmt"
	"log"
	"strings"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Group    string
	DBOption database.Option
}

func ListFeature(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	var features []database.FeatureConfig
	if option.Group == "" {
		features, err = database.ListFeatureConfig(db)
		if err != nil {
			log.Fatalf("failed listing feature configs: %v", err)
		}
	} else {
		features, err = database.ListFeatureConfigByGroup(db, option.Group)
		if err != nil {
			log.Fatalf("failed listing feature configs in group %s: %v", option.Group, err)
		}
	}

	fmt.Println(featureConfigTitle())
	for _, feature := range features {
		fmt.Println(feature.OneLineString())
	}
}

func featureConfigTitle() string {
	return strings.Join([]string{
		"Name", "Group", "Revision", "Status", "Category", "ValueType", "Description", "RevisionsLimit", "CreateTime", "ModifyTime"},
		",")
}
