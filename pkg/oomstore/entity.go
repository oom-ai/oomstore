package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) error {
	return s.metadata.CreateEntity(ctx, opt)
}

func (s *OomStore) GetEntity(ctx context.Context, name string) (*types.Entity, error) {
	return s.metadata.GetEntity(ctx, name)
}

func (s *OomStore) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	return s.metadata.ListEntity(ctx)
}

func (s *OomStore) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error {
	return s.metadata.UpdateEntity(ctx, opt)
}
