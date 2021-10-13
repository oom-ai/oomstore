package database

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateGroup(ctx context.Context, opt types.CreateGroupOpt) error {
	_, err := db.ExecContext(ctx,
		"insert into feature_group(name, entity_name, category, description) values(?, ?, ?, ?)",
		opt.Name, opt.EntityName, opt.Category, opt.Description)
	return err
}

func (db *DB) GetGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	var group types.FeatureGroup
	query := `SELECT * FROM feature_group WHERE name = ?`
	if err := db.GetContext(ctx, &group, query, groupName); err != nil {
		return nil, err
	}
	return &group, nil
}
