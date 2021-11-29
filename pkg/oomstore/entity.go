package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Get metadata of an entity by ID.
func (s *OomStore) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CacheGetEntity(ctx, id)
}

// Get metadata of an entity by name.
func (s *OomStore) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CacheGetEntityByName(ctx, name)
}

// List metadata of all entities.
func (s *OomStore) ListEntity(ctx context.Context) types.EntityList {
	_ = s.metadata.Refresh()
	return s.metadata.CacheListEntity(ctx)
}

// Create metadata for an entity.
func (s *OomStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (int, error) {
	if err := s.metadata.Refresh(); err != nil {
		return 0, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: opt,
	})
}

// Update metadata for an entity.
func (s *OomStore) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error {
	if err := s.metadata.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	entity, err := s.metadata.CacheGetEntityByName(ctx, opt.EntityName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateEntity(ctx, metadata.UpdateEntityOpt{
		EntityID:       entity.ID,
		NewDescription: opt.NewDescription,
	})
}
