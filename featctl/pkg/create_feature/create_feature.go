package create_feature

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"github.com/onestore-ai/onestore/pkg/database"
)

type Option struct {
	Name           string
	Group          string
	Category       string
	ValueType      string
	Revision       string
	RevisionsLimit int
	Status         string
	Description    string
	DBOption       database.Option
}

func createFeatureConfig(ctx context.Context, db *database.DB, option *Option) error {
	_, err := db.ExecContext(ctx,
		"insert into"+
			" feature_config(name, `group`, category, value_type, revision, revisions_limit, status, description)"+
			" values(?, ?, ?, ?, ?, ?, ?, ?)",
		option.Name,
		option.Group,
		option.Category,
		option.ValueType,
		option.Revision,
		option.RevisionsLimit,
		option.Status,
		option.Description,
	)
	return err
}

func getSourceTableName(ctx context.Context, db *database.DB, group string, revision string) (string, error) {
	var source string
	err := db.QueryRowContext(ctx,
		"select source from feature_revision where `group` = ? and revision = ?",
		group, revision).Scan(&source)
	if err == sql.ErrNoRows {
		return "", fmt.Errorf("revision not found: %s", revision)
	}
	return source, err
}

func getValueType(ctx context.Context, db *database.DB, option *Option) (string, error) {
	sourceTable, err := getSourceTableName(ctx, db, option.Group, option.Revision)
	if err != nil {
		return "", fmt.Errorf("failed fetching source table: %v", err)
	}

	column, err := db.ColumnInfo(ctx, sourceTable, option.Name)
	if err != nil {
		return "", fmt.Errorf("failed fetching source column: %v", err)
	}

	return column.Type, nil
}

func Create(ctx context.Context, option *Option) {
	db, err := database.Open(&option.DBOption)
	if err != nil {
		log.Fatalf("failed connecting feature store: %v", err)
	}
	defer db.Close()

	log.Println("obtainning value type...")
	valueType, err := getValueType(ctx, db, option)
	if err != nil {
		log.Fatalf("failed obtainning value type: %v", err)
	}
	option.ValueType = valueType

	log.Println("creating new feature...")
	if err = createFeatureConfig(ctx, db, option); err != nil {
		log.Fatalf("error creating new feature: %v", err)
	}

	log.Println("succeeded.")
}
