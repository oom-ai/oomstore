package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

// CreateEntity create an entity in the store
func (s *OneStore) CreateEntity(ctx context.Context, entityName, description string) (*types.Entity, error) {
	if err := s.db.CreateEntity(ctx, entityName, description); err != nil {
		return nil, err
	}

	return &types.Entity{
		Name:        entityName,
		Description: description,
	}, nil
}

// ListEntity: get all entity
func (s *OneStore) ListEntity(ctx context.Context) ([]*types.Entity, error) {
	entities, err := s.db.ListEntity(ctx)
	if err != nil {
		return nil, err
	}
	return entities, nil
}
