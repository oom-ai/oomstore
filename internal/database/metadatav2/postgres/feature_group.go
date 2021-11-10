package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateFeatureGroup(ctx context.Context, opt metadatav2.CreateFeatureGroupOpt) (int16, error) {
	if opt.Category != types.BatchFeatureCategory && opt.Category != types.StreamFeatureCategory {
		return 0, fmt.Errorf("illegal category %s, should be either 'stream' or 'batch'", opt.Category)
	}
	var featureGroupId int16
	query := "insert into feature_group(name, entity_id, category, description) values($1, $2, $3, $4) returning id"
	err := db.GetContext(ctx, &featureGroupId, query, opt.Name, opt.EntityID, opt.Category, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("feature group %s already exists", opt.Name)
			}
		}
	}
	return featureGroupId, err
}

func (db *DB) UpdateFeatureGroup(ctx context.Context, opt metadatav2.UpdateFeatureGroupOpt) error {
	and := make(map[string]interface{})
	if opt.Description != nil {
		and["description"] = *opt.Description
	}
	if opt.OnlineRevisionId != nil {
		and["online_revision_id"] = *opt.OnlineRevisionId
	}
	cond, args, err := dbutil.BuildConditions(and, nil)
	if err != nil {
		return err
	}
	args = append(args, opt.GroupID)

	if len(cond) == 0 {
		return fmt.Errorf("invalid option: nothing to update")
	}

	query := fmt.Sprintf("UPDATE feature_group SET %s WHERE id = ?", strings.Join(cond, ","))
	result, err := db.ExecContext(ctx, db.Rebind(query), args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("expect 1 affected row, get %d", rowsAffected)
	}
	return nil
}
