package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func createEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.CreateEntityOpt) (int, error) {
	var entityID int
	query := "insert into feature_entity(name, length, description) values($1, $2, $3) returning id"
	err := sqlxCtx.GetContext(ctx, &entityID, query, opt.EntityName, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("entity %s already exists", opt.EntityName)
		}
	}
	return entityID, err
}

func updateEntity(ctx context.Context, sqlxCtx metadata.SqlxContext, opt metadata.UpdateEntityOpt) error {
	query := "UPDATE feature_entity SET description = $1 WHERE id = $2"
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
