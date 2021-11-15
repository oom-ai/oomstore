package metadata

import (
	"context"

	"github.com/jmoiron/sqlx"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type StoreImpl struct {
	db *sqlx.DB

	tx *sqlx.Tx

	createEntity func(context.Context, ExtContext, CreateEntityOpt) (int16, error)
	updateEntity func(context.Context, ExtContext, UpdateEntityOpt) error
	getEntity    func(context.Context, int16) (*types.Entity, error)
}

func NewStoreImpl(db *sqlx.DB,
	createEntity func(context.Context, ExtContext, CreateEntityOpt) (int16, error),
	updateEntity func(context.Context, ExtContext, UpdateEntityOpt) error,
	getEntity func(context.Context, int16) (*types.Entity, error),

	// ...
) *StoreImpl {
	return &StoreImpl{
		db: db,

		// entity methods
		createEntity: createEntity,
		updateEntity: updateEntity,
		getEntity:    getEntity,

		// feature methods ...
	}
}

func (s *StoreImpl) CreateEntity(ctx context.Context, opt CreateEntityOpt) (int16, error) {
	if s.tx != nil {
		return s.createEntity(ctx, s.tx, opt)
	}
	return s.createEntity(ctx, s.db, opt)
}

func (s *StoreImpl) UpdateEntity(ctx context.Context, opt UpdateEntityOpt) error {
	if s.tx != nil {
		return s.updateEntity(ctx, s.tx, opt)
	}
	return s.updateEntity(ctx, s.db, opt)
}

func (s *StoreImpl) GetEntity(ctx context.Context, id int16) (*types.Entity, error) {
	return s.getEntity(ctx, id)
}

func (s *StoreImpl) copy() *StoreImpl {
	return &StoreImpl{
		createEntity: s.createEntity,
		updateEntity: s.updateEntity,
		getEntity:    s.getEntity,

		// ...
	}
}

func (s *StoreImpl) WithTransaction(ctx context.Context, fn func(context.Context, Store) error) (err error) {
	tx, err := s.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	txStore := s.copy()
	txStore.tx = tx

	defer func() {
		if p := recover(); p != nil {
			// a panic occurred, rollback and repanic
			_ = tx.Rollback()
			panic(p)
		} else if err != nil {
			// something went wrong, rollback
			_ = tx.Rollback()
		} else {
			// all good, commit
			err = tx.Commit()
		}
	}()

	return fn(ctx, txStore)
}
