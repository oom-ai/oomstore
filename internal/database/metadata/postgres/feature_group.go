package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
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

func (db *DB) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	query := "SELECT * FROM rich_feature_group WHERE name = $1"

	var group types.FeatureGroup
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
	query := "SELECT * FROM rich_feature_group"
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

func (db *DB) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error) {
	and := make(map[string]interface{})
	if opt.Description != nil {
		and["description"] = *opt.Description
	}
	if opt.OnlineRevisionId != nil {
		and["online_revision_id"] = *opt.OnlineRevisionId
	}
	cond, args, err := dbutil.BuildConditions(and, nil)
	if err != nil {
		return 0, err
	}
	args = append(args, opt.GroupName)

	if len(cond) == 0 {
		return 0, fmt.Errorf("invliad option: nothing to update")
	}

	query := fmt.Sprintf("UPDATE feature_group SET %s WHERE name = $%d", strings.Join(cond, ","), len(cond)+1)
	if result, err := db.ExecContext(ctx, db.Rebind(query), args...); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
