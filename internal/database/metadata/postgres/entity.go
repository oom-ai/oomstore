package postgres

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

func (db *DB) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	query := "insert into feature_entity(name, length, description) values($1, $2, $3)"
	_, err := db.ExecContext(ctx, query, opt.Name, opt.Length, opt.Description)
	return err
}

func (db *DB) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	var entity types.Entity
	query := "select * from feature_entity where name = $1"
	if err := db.GetContext(ctx, &entity, query, name); err != nil {
		return nil, err
	}
	return &entity, nil
}

func (db *DB) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	query := "select * from feature_entity"
	entities := make([]*types.Entity, 0)

	if err := db.SelectContext(ctx, &entities, query); err != nil {
		return nil, err
	}
	return entities, nil
}

func (db *DB) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error {
	query := "UPDATE feature_entity SET description = $1 WHERE name = $2"
	_, err := db.ExecContext(ctx, query, opt.NewDescription, opt.EntityName)
	return err
}
