package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt, category string) error {
	if category != types.BatchFeatureCategory && category != types.StreamFeatureCategory {
		return fmt.Errorf("illegal category %s, should be either 'stream' or 'batch'", category)
	}
	query := "insert into feature_group(name, entity_name, category, description) values($1, $2, $3, $4)"
	_, err := db.ExecContext(ctx, query, opt.Name, opt.EntityName, category, opt.Description)
	return err
}

func (db *DB) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	var group types.FeatureGroup
	query := `SELECT * FROM feature_group WHERE name = $1`
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
		query = query + " WHERE entity_name = $1"
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
	cond, args := buildUpdateFeatureGroupCond(opt)
	query := fmt.Sprintf("UPDATE feature_group SET %s WHERE name = $%d", strings.Join(cond, ","), len(cond)+1)
	_, err := db.ExecContext(ctx, query, args...)
	return err
}

func buildUpdateFeatureGroupCond(opt types.UpdateFeatureGroupOpt) ([]string, []interface{}) {
	cond := make([]string, 0)
	args := make([]interface{}, 0)
	var id int
	if opt.Description != nil {
		id++
		cond = append(cond, fmt.Sprintf("description = $%d", id))
		args = append(args, *opt.Description)
	}
	if opt.Revision != nil {
		id++
		cond = append(cond, fmt.Sprintf("revision = $%d", id))
		args = append(args, *opt.Revision)
	}
	if opt.DataTable != nil {
		id++
		cond = append(cond, fmt.Sprintf("data_table = $%d", id))
		args = append(args, *opt.DataTable)
	}
	args = append(args, opt.GroupName)
	return cond, args
}

func (db *DB) UpdateFeatureGroupRevision(ctx context.Context, revision int64, dataTable string, groupName string) error {
	cmd := "UPDATE feature_group SET revision = $1, data_table = $2 WHERE name = $3"
	_, err := db.ExecContext(ctx, cmd, revision, dataTable, groupName)
	return err
}
