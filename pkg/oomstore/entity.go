package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// CreateEntity create an entity in the store
func (s *OomStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (*types.Entity, error) {
	if err := s.metadata.CreateEntity(ctx, opt); err != nil {
		return nil, err
	}
	return s.GetEntity(ctx, opt.Name)
}

func (s *OomStore) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return s.metadata.GetEntity(ctx, name)
}

// ListEntity: get all entity
func (s *OomStore) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return s.metadata.ListEntity(ctx)
}

func (s *OomStore) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error {
	return s.metadata.UpdateEntity(ctx, opt)
}
