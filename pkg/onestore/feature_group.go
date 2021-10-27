package onestore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

func (s *OneStore) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt) (*types.FeatureGroup, error) {
	if err := s.metadata.CreateFeatureGroup(ctx, metadata.CreateFeatureGroupOpt{
		CreateFeatureGroupOpt: opt,
		Category:              types.BatchFeatureCategory,
	}); err != nil {
		return nil, err
	}
	return s.GetFeatureGroup(ctx, opt.Name)
}

func (s *OneStore) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	return s.metadata.GetFeatureGroup(ctx, groupName)
}

func (s *OneStore) ListFeatureGroup(ctx context.Context, entityName *string) ([]*types.FeatureGroup, error) {
	return s.metadata.ListFeatureGroup(ctx, entityName)

}

func (s *OneStore) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) error {
	return s.metadata.UpdateFeatureGroup(ctx, opt)
}
