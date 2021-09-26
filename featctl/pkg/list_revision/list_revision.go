package list_revision

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

func ListRevision(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	revisions, err := database.ListGroupRevisionByGroup(db, option.Group)
	if err != nil {
		log.Fatalf("failed listing revisions for group %s: %v", option.Group, err)
	}

	fmt.Println(groupRevisionTitle())
	for _, revision := range revisions {
		fmt.Println(revision.OneLineString())
	}
}

func groupRevisionTitle() string {
	return strings.Join([]string{"Group", "Revision", "Source", "Description", "CreateTime", "ModifyTime"},
		",")
}
