package database

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) CreateEntity(ctx context.Context, entityName, description string) error {
	_, err := db.ExecContext(ctx,
		"insert into feature_entity(name, description) values(?, ?)",
		entityName, description)
	return err
}

func (db *DB) ListEntity(ctx context.Context) ([]types.Entity, error) {
	query := "select name, description, create_time, modify_time from feature_entity"
	entities := make([]types.Entity, 0)

	if err := db.SelectContext(ctx, &entities, query); err != nil {
		return nil, err
	}
	return entities, nil
}
