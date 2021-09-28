package update_feature

import (
	"context"
	"database/sql"
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
			_, err := db.ExecContext(ctx,
				fmt.Sprintf("update feature_config set %s = ?", item.field)+
					" where `group` = ? and name = ?", item.value, option.Group, option.Name)
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
		return requireRevisionExists(ctx, db, option.Group, option.Revision)
	}
	return nil
}

func requireRevisionExists(ctx context.Context, db *database.DB, group string, revision string) error {
	err := db.GetContext(ctx, &revision,
		"select revision from feature_revision where `group` = ? and revision = ?",
		group, revision)
	if err == sql.ErrNoRows {
		return fmt.Errorf("revision '%s' not found int feature group '%s'", revision, group)
	}
	return err
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
