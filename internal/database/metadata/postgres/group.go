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
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateGroupOpt) (int, error) {
	if opt.Category != types.BatchFeatureCategory && opt.Category != types.StreamFeatureCategory {
		return 0, fmt.Errorf("illegal category '%s', should be either 'stream' or 'batch'", opt.Category)
	}
	var groupID int
	query := "insert into feature_group(name, entity_id, category, description) values($1, $2, $3, $4) returning id"
	err := sqlxCtx.GetContext(ctx, &groupID, query, opt.GroupName, opt.EntityID, opt.Category, opt.Description)
	if err != nil {
		if e2, ok := err.(*pq.Error); ok {
			if e2.Code == pgerrcode.UniqueViolation {
				return 0, fmt.Errorf("feature group %s already exists", opt.GroupName)
			}
		}
	}
	return groupID, err
}

func updateGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateGroupOpt) error {
	and := make(map[string]interface{})
	if opt.NewDescription != nil {
		and["description"] = *opt.NewDescription
	}
	if opt.NewOnlineRevisionID != nil {
		and["online_revision_id"] = *opt.NewOnlineRevisionID
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
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), args...)
	if err != nil {
		return err
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}
	if rowsAffected != 1 {
		return fmt.Errorf("failed to update feature group %d: feature group not found", opt.GroupID)
	}
	return nil
}

func getGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Group, error) {
	var group types.Group
	query := `SELECT * FROM feature_group WHERE id = $1`
	if err := sqlxCtx.GetContext(ctx, &group, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature group %d not found", id))
		}
		return nil, err
	}

	entity, err := getEntity(ctx, sqlxCtx, group.EntityID)
	if err != nil {
		return nil, err
	}
	group.Entity = entity
	return &group, nil
}

func getGroupByName(ctx context.Context, sqlxCtx metadata.SqlxContext, name string) (*types.Group, error) {
	var group types.Group
	query := `SELECT * FROM feature_group WHERE name = $1`
	if err := sqlxCtx.GetContext(ctx, &group, query, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature group %s not found", name))
		}
		return nil, err
	}

	entity, err := getEntity(ctx, sqlxCtx, group.EntityID)
	if err != nil {
		return nil, err
	}
	group.Entity = entity
	return &group, nil
}

func listGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, entityID *int) (types.GroupList, error) {
	query := "SELECT * FROM feature_group"
	args := make([]interface{}, 0)
	if entityID != nil {
		query = fmt.Sprintf("%s WHERE entity_id = $1", query)
		args = append(args, *entityID)
	}
	var groups types.GroupList
	if err := sqlxCtx.SelectContext(ctx, &groups, query, args...); err != nil {
		return nil, err
	}

	entityIDs := groups.EntityIDs()
	entities, err := listEntity(ctx, sqlxCtx, &entityIDs)
	if err != nil {
		return nil, err
	}
	for _, group := range groups {
		entity := entities.Find(func(e *types.Entity) bool {
			return group.EntityID == e.ID
		})
		if entity == nil {
			return nil, fmt.Errorf("cannot find entity %d for group %d", group.EntityID, group.ID)
		}
		group.Entity = entity
	}
	return groups, nil
}
