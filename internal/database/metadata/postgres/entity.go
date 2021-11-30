package postgres

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	var entityID int
	query := "insert into entity(name, length, description) values($1, $2, $3) returning id"
	err := sqlxCtx.GetContext(ctx, &entityID, query, opt.EntityName, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("entity %s already exists", opt.EntityName)
		}
	}
	return entityID, err
}

func updateEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateEntityOpt) error {
	if opt.NewDescription == nil {
		return fmt.Errorf("invalid option: nothing to update")
	}

	query := "UPDATE entity SET description = $1 WHERE id = $2"
	result, err := sqlxCtx.ExecContext(ctx, query, opt.NewDescription, opt.EntityID)
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

func getEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, id int) (*types.Entity, error) {
	var entity types.Entity
	query := "select * from entity where id = $1"
	if err := sqlxCtx.GetContext(ctx, &entity, query, id); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature entity %d not found", id))
		}
		return nil, err
	}

	return &entity, nil
}

func getEntityByName(ctx context.Context, sqlxCtx metadata.SqlxContext, name string) (*types.Entity, error) {
	var entity types.Entity
	query := "select * from entity where name = $1"
	if err := sqlxCtx.GetContext(ctx, &entity, query, name); err != nil {
		if err == sql.ErrNoRows {
			return nil, errdefs.NotFound(fmt.Errorf("feature entity %s not found", name))
		}
		return nil, err
	}

	return &entity, nil

}

func listEntity(ctx context.Context, sqlxCtx metadata.SqlxContext) (types.EntityList, error) {
	query := "select * from entity"
	entities := types.EntityList{}

	if err := sqlxCtx.SelectContext(ctx, &entities, query); err != nil {
		return nil, err
	}
	return entities, nil
}
