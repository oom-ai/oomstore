package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt) (*types.FeatureGroup, error) {
	if err := s.db.CreateFeatureGroup(ctx, opt, types.BatchFeatureCategory); err != nil {
		return nil, err
	}
	return s.GetFeatureGroup(ctx, opt.Name)
}

func (s *OneStore) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	return s.db.GetFeatureGroup(ctx, groupName)
}
