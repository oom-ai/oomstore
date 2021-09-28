package update_feature

import (
	"context"
	"fmt"
	"log"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Name  string
	Group string

	Revision        string
	RevisionChanged bool

	RevisionsLimit        int
	RevisionsLimitChanged bool

	Status        string
	StatusChanged bool

	Description        string
	DescriptionChanged bool

	DBOption database.Option
}

func updateFeature(ctx context.Context, db *database.DB, option *Option) error {
	updated := 0
	for _, item := range []struct {
		field     string
		value     interface{}
		condition bool
	}{
		{"revision", option.Revision, option.RevisionChanged},
		{"revisions_limit", option.RevisionsLimit, option.RevisionsLimitChanged},
		{"status", option.Status, option.StatusChanged},
		{"description", option.Description, option.DescriptionChanged},
	} {
		if item.condition {
			_, err := database.UpdateFeatureConfig(ctx, db, item.field, item.value,
				// where clause
				option.Group,
				option.Name,
			)
			if err != nil {
				return err
			}
			updated++
		}
	}

	if updated == 0 {
		return fmt.Errorf("nothing to set\n")
	}

	return nil
}

func validateOptions(ctx context.Context, db *database.DB, option *Option) error {
	if option.RevisionChanged {
		return database.RevisionExists(ctx, db, option.Group, option.Revision)
	}
	return nil
}

func Run(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	if err = validateOptions(ctx, db, option); err != nil {
		log.Fatalf("failed validating options: %v", err)
	}

	log.Println("updating feature...")
	if err = updateFeature(ctx, db, option); err != nil {
		log.Fatalf("failed updating feature: %v", err)
	}

	log.Println("succeeded.")
}
