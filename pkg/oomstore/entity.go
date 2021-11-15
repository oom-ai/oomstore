package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) GetEntity(ctx context.Context, id int16) (*types.Entity, error) {
	return s.metadata.GetEntity(ctx, id)
}

func (s *OomStore) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return s.metadata.GetEntityByName(ctx, name)
}

func (s *OomStore) ListEntity(ctx context.Context) types.EntityList {
	return s.metadata.ListEntity(ctx)
}

func (s *OomStore) CreateEntity(ctx context.Context, opt metadata.CreateEntityOpt) (int16, error) {
	return s.metadata.CreateEntity(ctx, opt)
}

func (s *OomStore) UpdateEntity(ctx context.Context, opt metadata.UpdateEntityOpt) error {
	return s.metadata.UpdateEntity(ctx, opt)
}
