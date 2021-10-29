package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) error {
	if opt.Category != types.BatchFeatureCategory && opt.Category != types.StreamFeatureCategory {
		return fmt.Errorf("illegal category %s, should be either 'stream' or 'batch'", opt.Category)
	}
	query := "insert into feature_group(name, entity_name, category, description) values($1, $2, $3, $4)"
	_, err := db.ExecContext(ctx, query, opt.Name, opt.EntityName, opt.Category, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return fmt.Errorf("feature group %s already exist", opt.Name)
			}
		}
	}
	return err
}

func (db *DB) getFeatureGroup(ctx context.Context, groupName string, source string, group interface{}) error {
	query := fmt.Sprintf(`SELECT * FROM "%s" WHERE name = $1`, source)
	if err := db.GetContext(ctx, group, query, groupName); err != nil {
		if err == sql.ErrNoRows {
			return fmt.Errorf("group %s does not exist", groupName)
		}
		return err
	}
	return nil
}

func (db *DB) listFeatureGroup(ctx context.Context, entityName *string, source string, groups interface{}) error {
	var cond []interface{}
	query := fmt.Sprintf("SELECT * FROM %s", source)
	if entityName != nil {
		query = query + " WHERE entity_name = $1"
		cond = append(cond, *entityName)
	}

	if err := db.SelectContext(ctx, groups, query, cond...); err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		return err
	}
	return nil
}

func (db *DB) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	var group types.FeatureGroup
	if err := db.getFeatureGroup(ctx, groupName, "feature_group", &group); err != nil {
		return nil, err
	}
	return &group, nil
}

func (db *DB) ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error) {
	var groups []*types.FeatureGroup
	return groups, db.listFeatureGroup(ctx, entityName, "feature_group", &groups)
}

func (db *DB) GetRichFeatureGroup(ctx context.Context, groupName string) (*types.RichFeatureGroup, error) {
	var group types.RichFeatureGroup
	return &group, db.getFeatureGroup(ctx, groupName, "rich_feature_group", &group)

}

func (db *DB) ListRichFeatureGroup(ctx context.Context, entityName *string) ([]*types.RichFeatureGroup, error) {
	var groups []*types.RichFeatureGroup
	return groups, db.listFeatureGroup(ctx, entityName, "rich_feature_group", &groups)
}

func (db *DB) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error) {
	cond, args := buildUpdateFeatureGroupCond(opt)
	if len(cond) == 0 {
		return 0, fmt.Errorf("invliad option: nothing to update")
	}

	query := fmt.Sprintf("UPDATE feature_group SET %s WHERE name = $%d", strings.Join(cond, ","), len(cond)+1)
	if result, err := db.ExecContext(ctx, query, args...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
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
	if opt.OnlineRevisionId != nil {
		id++
		cond = append(cond, fmt.Sprintf("online_revision_id = $%d", id))
		args = append(args, *opt.OnlineRevisionId)
	}
	args = append(args, opt.GroupName)
	return cond, args
}
