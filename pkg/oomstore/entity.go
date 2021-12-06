package oomstore

import (
	"context"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

// Get metadata of an entity by ID.
func (s *OomStore) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	return s.metadata.GetEntity(ctx, id)
}

// Get metadata of an entity by name.
func (s *OomStore) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return s.metadata.GetEntityByName(ctx, name)
}

// List metadata of all entities.
func (s *OomStore) ListEntity(ctx context.Context) (types.EntityList, error) {
	return s.metadata.ListEntity(ctx, nil)
}

// Create metadata for an entity.
func (s *OomStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (int, error) {
	return s.metadata.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: opt,
	})
}

// Update metadata for an entity.
func (s *OomStore) UpdateEntity(ctx context.Context, opt types.UpdateEntityOpt) error {
	entity, err := s.metadata.GetEntityByName(ctx, opt.EntityName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateEntity(ctx, metadata.UpdateEntityOpt{
		EntityID:       entity.ID,
		NewDescription: opt.NewDescription,
	})
}
