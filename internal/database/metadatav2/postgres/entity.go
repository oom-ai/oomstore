package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (*types.Entity, error) {
	var entity types.Entity
	query := "insert into feature_entity(name, length, description) values($1, $2, $3) returning *"
	err := db.GetContext(ctx, &entity, query, opt.Name, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return nil, fmt.Errorf("entity %s already exists", opt.Name)
		}
	}
	return &entity, err
}

func (db *DB) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error) {
	query := "UPDATE feature_entity SET description = $1 WHERE name = $2"
	if result, err := db.ExecContext(ctx, query, opt.NewDescription, opt.EntityName); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
