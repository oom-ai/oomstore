package sqlutil

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func UpdateEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateEntityOpt) error {
	if opt.NewDescription == nil {
		return errors.Errorf("invalid option: nothing to update")
	}

	query := "UPDATE entity SET description = ? WHERE id = ?"
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), opt.NewDescription, opt.EntityID)
	if err != nil {
		return errors.WithStack(err)
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return errors.WithStack(err)
	}

	if rowsAffected != 1 {
		return errors.Errorf("failed to update entity %d: entity not found", opt.EntityID)
	}
	return nil
}

func GetEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Entity, error) {
	var entity types.Entity
	query := "SELECT * FROM entity WHERE id = ?"
	if err := sqlxCtx.GetContext(ctx, &entity, sqlxCtx.Rebind(query), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errors.Errorf("feature entity %d not found", id))
		}
		return nil, errors.WithStack(err)
	}

	return &entity, nil
}

func GetEntityByName(ctx context.Context, sqlxCtx metadata.SqlxContext, name string) (*types.Entity, error) {
	var entity types.Entity
	query := "SELECT * FROM entity WHERE name = ?"
	if err := sqlxCtx.GetContext(ctx, &entity, sqlxCtx.Rebind(query), name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(errors.Errorf("feature entity %s not found", name))
		}
		return nil, errors.WithStack(err)
	}

	return &entity, nil
}

func ListEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, entityIDs *[]int) (types.EntityList, error) {
	query := "SELECT * FROM entity"
	var args []interface{}
	var err error
	if entityIDs != nil {
		if len(*entityIDs) == 0 {
			return nil, nil
		}
		query, args, err = sqlx.In(fmt.Sprintf("%s WHERE id IN (?)", query), *entityIDs)
		if err != nil {
			return nil, errors.WithStack(err)
		}
	}
	entities := types.EntityList{}
	if err := sqlxCtx.SelectContext(ctx, &entities, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, errors.WithStack(err)
	}
	return entities, nil
}
