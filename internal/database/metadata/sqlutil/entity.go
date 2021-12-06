package sqlutil

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/pkg/errdefs"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/jmoiron/sqlx"
)

func UpdateEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateEntityOpt) error {
	if opt.NewDescription == nil {
		return fmt.Errorf("invalid option: nothing to update")
	}

	query := "UPDATE entity SET description = ? WHERE id = ?"
	result, err := sqlxCtx.ExecContext(ctx, sqlxCtx.Rebind(query), opt.NewDescription, opt.EntityID)
	if err != nil {
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil {
		return err
	}

	if rowsAffected != 1 {
		return fmt.Errorf("failed to update entity %d: entity not found", opt.EntityID)
	}
	return nil
}

func GetEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Entity, error) {
	var entity types.Entity
	query := "SELECT * FROM entity WHERE id = ?"
	if err := sqlxCtx.GetContext(ctx, &entity, sqlxCtx.Rebind(query), id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature entity %d not found", id))
		}
		return nil, err
	}

	return &entity, nil
}

func GetEntityByName(ctx context.Context, sqlxCtx metadata.SqlxContext, name string) (*types.Entity, error) {
	var entity types.Entity
	query := "SELECT * FROM entity WHERE name = ?"
	if err := sqlxCtx.GetContext(ctx, &entity, sqlxCtx.Rebind(query), name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature entity %s not found", name))
		}
		return nil, err
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
			return nil, err
		}
	}
	entities := types.EntityList{}
	if err := sqlxCtx.SelectContext(ctx, &entities, sqlxCtx.Rebind(query), args...); err != nil {
		return nil, err
	}
	return entities, nil
}
