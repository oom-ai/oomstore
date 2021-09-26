package list_feature

import (
	"context"
	"fmt"
	"log"
	"strings"

	database2 "github.com/onestore-ai/onestore/featctl/pkg/database"
	"github.com/onestore-ai/onestore/featctl/pkg/utils"
	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Group    string
	DBOption database2.Option
}

func ListFeature(ctx context.Context, option *Option) {
	sqlxDBOption := utils.BuildSqlxDBOption(option.DBOption)
	db, err := database.Open(sqlxDBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	features, err := database.ListFeatureConfigByGroup(db, option.Group)
	if err != nil {
		log.Fatalf("failed listing feature configs in group %s: %v", option.Group, err)
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
