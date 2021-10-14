package database

import (
	"context"
	"fmt"
	"strings"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateFeature(ctx context.Context, opt types.CreateFeatureOpt) error {
	query := "INSERT INTO feature(name, group_name, value_type, description) VALUES (?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query, opt.FeatureName, opt.GroupName, opt.ValueType, opt.Description)
	return err
}

func (db *DB) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	var feature types.Feature
	query := `SELECT * FROM feature WHERE name = ?`
	if err := db.GetContext(ctx, &feature, query, featureName); err != nil {
		return nil, err
	}
	return &feature, nil
}

func (db *DB) ListFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.Feature, error) {
	query := "SELECT * FROM feature"
	cond, args := buildListFeatureCond(opt)
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}

	features := make([]*types.Feature, 0)
	if err := db.SelectContext(ctx, &features, query, args); err != nil {
		return nil, err
	}
	return features, nil
}

func (db *DB) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	query := "UPDATE feature SET description = ? WHERE name = ?"
	_, err := db.ExecContext(ctx, query, opt.NewDescription, opt.FeatureName)
	return err
}

func buildListFeatureCond(opt types.ListFeatureOpt) (string, []string) {
	cond := make([]string, 0)
	args := make([]string, 0)
	if opt.EntityName != nil {
		cond = append(cond, "entity_name = ?")
		args = append(args, *opt.EntityName)
	}
	if opt.GroupName != nil {
		cond = append(cond, "group_name = ?")
		args = append(args, *opt.GroupName)
	}
	return strings.Join(cond, " AND "), args
}
