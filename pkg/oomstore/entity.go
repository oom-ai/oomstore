package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// GetEntity gets metadata of an entity by ID.
func (s *OomStore) GetEntity(ctx context.Context, id int) (*types.Entity, error) {
	return s.metadata.GetEntity(ctx, id)
}

// GetEntityByName gets metadata of an entity by name.
func (s *OomStore) GetEntityByName(ctx context.Context, name string) (*types.Entity, error) {
	return s.metadata.GetEntityByName(ctx, name)
}

// ListEntity lists metadata of entities given entity names.
func (s *OomStore) ListEntity(ctx context.Context, opt types.ListEntityOpt) (types.EntityList, error) {
	entities, err := s.metadata.ListEntity(ctx, nil)
	if err != nil {
		return nil, err
	}
	if opt.EntityNames != nil {
		entities = entities.Filter(func(e *types.Entity) bool {
			for _, name := range *opt.EntityNames {
				if e.Name == name {
					return true
				}
			}
			return false
		})
	}
	return entities, nil
}

// CreateEntity creates metadata for an entity.
func (s *OomStore) CreateEntity(ctx context.Context, opt types.CreateEntityOpt) (int, error) {
	return s.metadata.CreateEntity(ctx, metadata.CreateEntityOpt{
		CreateEntityOpt: opt,
	})
}

// UpdateEntity updates metadata for an entity.
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
