package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt) error {
	entity, err := s.metadata.GetEntity(ctx, opt.EntityName)
	if err != nil {
		return err
	}
	return s.metadata.CreateFeatureGroup(ctx, metadata.CreateFeatureGroupOpt{
		Name:        opt.Name,
		EntityId:    entity.ID,
		Description: opt.Description,
		Category:    types.BatchFeatureCategory,
	})
}

func (s *OomStore) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	return s.metadata.GetFeatureGroup(ctx, groupName)
}

func (s *OomStore) ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error) {
	return s.metadata.ListFeatureGroup(ctx, entityName)
}

func (s *OomStore) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) (int64, error) {
	return s.metadata.UpdateFeatureGroup(ctx, opt)
}
