package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (int16, error) {
	var entityId int16
	query := "insert into feature_entity(name, length, description) values($1, $2, $3) returning id"
	err := db.GetContext(ctx, &entityId, query, opt.Name, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return 0, fmt.Errorf("entity %s already exists", opt.Name)
		}
	}
	return entityId, err
}

func (db *DB) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error) {
	query := "UPDATE feature_entity SET description = $1 WHERE name = $2"
	if result, err := db.ExecContext(ctx, query, opt.NewDescription, opt.EntityName); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
