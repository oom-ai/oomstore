package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func UpdateGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateGroupOpt) error {
	and := make(map[string]interface{})
	if opt.NewDescription != nil {
		and["description"] = *opt.NewDescription
	}
	if opt.NewSnapshotInterval != nil {
		and["snapshot_interval"] = *opt.NewSnapshotInterval
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
		return errdefs.Errorf("invalid option: nothing to update")
	}

	query := fmt.Sprintf("UPDATE feature_group SET %s WHERE id = ?", strings.Join(cond, ","))
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), args...)
	if err != nil {
		return errdefs.WithStack(err)
	}
	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errdefs.WithStack(err)
	}
	if rowsAffected != 1 {
		return errdefs.Errorf("failed to update feature group %d: feature group not found", opt.GroupID)
	}
	return nil
}

func GetGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Group, error) {
	var group types.Group
	query := `SELECT * FROM feature_group WHERE id = ?`
	if err := sqlxCtx.GetContext(ctx, &group, sqlxCtx.Rebind(query), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errdefs.Errorf("feature group %d not found", id))
		}
		return nil, errdefs.WithStack(err)
	}

	entity, err := GetEntity(ctx, sqlxCtx, group.EntityID)
	if err != nil {
		return nil, err
	}
	group.Entity = entity
	return &group, nil
}

func GetGroupByName(ctx context.Context, sqlxCtx metadata.SqlxContext, name string) (*types.Group, error) {
	var group types.Group
	query := `SELECT * FROM feature_group WHERE name = ?`
	if err := sqlxCtx.GetContext(ctx, &group, sqlxCtx.Rebind(query), name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errdefs.Errorf("feature group %s not found", name))
		}
		return nil, errdefs.WithStack(err)
	}

	entity, err := GetEntity(ctx, sqlxCtx, group.EntityID)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	group.Entity = entity
	return &group, nil
}

func ListGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.ListGroupOpt) (types.GroupList, error) {
	cond, args, err := buildListGroupCond(opt)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}

	query := `SELECT * FROM feature_group`
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}
	query = fmt.Sprintf("%s ORDER BY id ASC", query)
	var groups types.GroupList
	if err := sqlxCtx.SelectContext(ctx, &groups, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, errdefs.WithStack(err)
	}

	if err := enrichGroups(ctx, sqlxCtx, groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func buildListGroupCond(opt metadata.ListGroupOpt) (cond string, args []interface{}, err error) {
	conds := make([]string, 0)

	if opt.EntityIDs != nil {
		if len(*opt.EntityIDs) == 0 {
			return "false", args, nil
		}
		s, inArgs, err := sqlx.In("entity_id IN (?)", *opt.EntityIDs)
		if err != nil {
			return "", nil, errdefs.WithStack(err)
		}
		conds = append(conds, s)
		args = append(args, inArgs...)
	}
	if opt.GroupIDs != nil {
		if len(*opt.GroupIDs) == 0 {
			return "false", args, nil
		}
		s, inArgs, err := sqlx.In("id IN (?)", *opt.GroupIDs)
		if err != nil {
			return "", nil, errdefs.WithStack(err)
		}
		conds = append(conds, s)
		args = append(args, inArgs...)
	}
	return strings.Join(conds, " AND "), args, nil
}

func enrichGroups(ctx context.Context, sqlxCtx metadata.SqlxContext, groups types.GroupList) error {
	entityIDs := groups.EntityIDs()
	entities, err := ListEntity(ctx, sqlxCtx, &entityIDs)
	if err != nil {
		return errdefs.WithStack(err)
	}
	for _, group := range groups {
		entity := entities.Find(func(e *types.Entity) bool {
			return group.EntityID == e.ID
		})
		if entity == nil {
			return errdefs.Errorf("cannot find entity %d for group %d", group.EntityID, group.ID)
		}
		group.Entity = entity
	}
	return nil
}
