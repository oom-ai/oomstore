package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgerrcode"
	"github.com/lib/pq"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	return createEntity(ctx, db, opt)
}

func (tx *Tx) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	return createEntity(ctx, tx, opt)
}

func createEntity(ctx context.Context, ext metadata.ExtContext, opt types.CreateEntityOpt) error {
	query := "insert into feature_entity(name, length, description) values($1, $2, $3)"
	_, err := ext.ExecContext(ctx, query, opt.Name, opt.Length, opt.Description)
	if er, ok := err.(*pq.Error); ok {
		if er.Code == pgerrcode.UniqueViolation {
			return fmt.Errorf("entity %s already exists", opt.Name)
		}
	}
	return err
}

func (db *DB) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return getEntity(ctx, db, name)
}

func (tx *Tx) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return getEntity(ctx, tx, name)
}

func getEntity(ctx context.Context, ext metadata.ExtContext, name string) (*types.Entity, error) {
	var entity types.Entity
	query := "select * from feature_entity where name = $1"
	if err := ext.GetContext(ctx, &entity, query, name); err != nil {
		return nil, err
	}

	return &entity, nil
}

func (db *DB) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return listEntity(ctx, db)
}

func (tx *Tx) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return listEntity(ctx, tx)
}

func listEntity(ctx context.Context, ext metadata.ExtContext) ([]*types.Entity, error) {
	query := "select * from feature_entity"
	entities := make([]*types.Entity, 0)

	if err := ext.SelectContext(ctx, &entities, query); err != nil {
		return nil, err
	}
	return entities, nil
}

func (db *DB) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error) {
	return updateEntity(ctx, db, opt)
}

func (tx *Tx) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) (int64, error) {
	return updateEntity(ctx, tx, opt)
}

func updateEntity(ctx context.Context, ext metadata.ExtContext, opt types.UpdateEntityOpt) (int64, error) {
	query := "UPDATE feature_entity SET description = $1 WHERE name = $2"
	if result, err := ext.ExecContext(ctx, query, opt.NewDescription, opt.EntityName); err != nil {
		return 0, err
	} else {
		return result.RowsAffected()
	}
}
