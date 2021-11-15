package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
)

func CreateEntity(ctx context.Context, ext metadata.ExtContext, opt metadata.CreateEntityOpt) (int16, error) {
	var entityId int16
	query := "insert into feature_entity(name, length, description) values($1, $2, $3) returning id"
	err := ext.GetContext(ctx, &entityId, query, opt.Name, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("entity %s already exists", opt.Name)
		}
	}
	return entityId, err
}

func UpdateEntity(ctx context.Context, ext metadata.ExtContext, opt metadata.UpdateEntityOpt) error {
	query := "UPDATE feature_entity SET description = $1 WHERE id = $2"
	result, err := ext.ExecContext(ctx, query, opt.NewDescription, opt.EntityID)
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
