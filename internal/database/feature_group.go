package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt, category string) error {
	if category != types.BatchFeatureCategory && category != types.StreamFeatureCategory {
		return fmt.Errorf("illegal category %s, should be either 'stream' or 'batch'", category)
	}
	query := "insert into feature_group(name, entity_name, category, description) values(?, ?, ?, ?)"
	_, err := db.ExecContext(ctx, query, opt.Name, opt.EntityName, category, opt.Description)
	return err
}

func (db *DB) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	var group types.FeatureGroup
	query := `SELECT * FROM feature_group WHERE name = ?`
	if err := db.GetContext(ctx, &group, query, groupName); err != nil {
		if err == sql.ErrNoRows {
			return nil, fmt.Errorf("group %s does not exist", groupName)
		}
		return nil, err
	}
	return &group, nil
}

func (db *DB) ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error) {
	var cond []interface{}
	query := "SELECT * FROM feature_group"
	if entityName != nil {
		query = query + " WHERE entity_name = ?"
		cond = append(cond, *entityName)
	}

	var groups []*types.FeatureGroup
	if err := db.SelectContext(ctx, &groups, query, cond...); err != nil {
		if err == sql.ErrNoRows {
			return nil, nil
		}
		return nil, err
	}
	return groups, nil

}

func (db *DB) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) error {
	query := "UPDATE feature_group SET description = ? WHERE name = ?"
	_, err := db.ExecContext(ctx, query, opt.NewDescription, opt.GroupName)
	return err
}

func UpdateFeatureGroupRevision(ctx context.Context, tx *sqlx.Tx, revision int64, dataTable string, groupName string) error {
	cmd := "UPDATE feature_group SET revision = ?, data_table = ? WHERE name = ?"
	_, err := tx.ExecContext(ctx, cmd, revision, dataTable, groupName)
	return err
}
