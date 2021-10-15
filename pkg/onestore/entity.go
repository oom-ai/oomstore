package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

// CreateEntity create an entity in the store
func (s *OneStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (*types.Entity, error) {
	if err := s.db.CreateEntity(ctx, opt); err != nil {
		return nil, err
	}
	return s.GetEntity(ctx, opt.Name)
}

func (s *OneStore) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return s.db.GetEntity(ctx, name)
}

// ListEntity: get all entity
func (s *OneStore) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return s.db.ListEntity(ctx)
}
