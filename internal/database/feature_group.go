package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jmoiron/sqlx"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt, category string) error {
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

func UpdateFeatureGroup(ctx context.Context, tx *sqlx.Tx, groupName string, revision int64, dataTable string) error {
	cmd := "UPDATE feature_group SET revision = ? AND data_table = ? WHERE name = ?"
	_, err := tx.ExecContext(ctx, cmd, revision, dataTable, groupName)
	return err
}
