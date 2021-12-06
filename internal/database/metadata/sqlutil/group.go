package sqlutil

import (
	"context"
	"database/sql"
	"fmt"
	"strings"

	"github.com/ethhte88/oomstore/internal/database/dbutil"
	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/pkg/errdefs"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func UpdateGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateGroupOpt) error {
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

func GetGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Group, error) {
	var group types.Group
	query := `SELECT * FROM feature_group WHERE id = ?`
	if err := sqlxCtx.GetContext(ctx, &group, sqlxCtx.Rebind(query), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature group %d not found", id))
		}
		return nil, err
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
			return nil, errdefs.NotFound(fmt.Errorf("feature group %s not found", name))
		}
		return nil, err
	}

	entity, err := GetEntity(ctx, sqlxCtx, group.EntityID)
	if err != nil {
		return nil, err
	}
	group.Entity = entity
	return &group, nil
}

func ListGroup(ctx context.Context, sqlxCtx metadata.SqlxContext, entityID *int, groupIDs *[]int) (types.GroupList, error) {
	cond, args, err := buildListGroupCond(entityID, groupIDs)
	if err != nil {
		return nil, err
	}

	query := `SELECT * FROM feature_group`
	if len(cond) > 0 {
		query = fmt.Sprintf("%s WHERE %s", query, cond)
	}
	var groups types.GroupList
	if err := sqlxCtx.SelectContext(ctx, &groups, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, err
	}

	if err := enrichGroups(ctx, sqlxCtx, groups); err != nil {
		return nil, err
	}
	return groups, nil
}

func buildListGroupCond(entityID *int, groupIDs *[]int) (string, []interface{}, error) {
	args := make([]interface{}, 0)
	cond := make([]string, 0)

	if entityID != nil {
		cond = append(cond, "entity_id = ?")
		args = append(args, *entityID)
	}
	if groupIDs != nil {
		if len(*groupIDs) == 0 {
			return "false", args, nil
		}
		s, inArgs, err := sqlx.In("id IN (?)", *groupIDs)
		if err != nil {
			return "", nil, err
		}
		cond = append(cond, s)
		args = append(args, inArgs...)
	}
	return strings.Join(cond, " AND "), args, nil
}

func enrichGroups(ctx context.Context, sqlxCtx metadata.SqlxContext, groups types.GroupList) error {
	entityIDs := groups.EntityIDs()
	entities, err := ListEntity(ctx, sqlxCtx, &entityIDs)
	if err != nil {
		return err
	}
	for _, group := range groups {
		entity := entities.Find(func(e *types.Entity) bool {
			return group.EntityID == e.ID
		})
		if entity == nil {
			return fmt.Errorf("cannot find entity %d for group %d", group.EntityID, group.ID)
		}
		group.Entity = entity
	}
	return nil
}
