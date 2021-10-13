package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

// GetFeature: get feature by featureName
func (s *OneStore) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	feature, err := s.db.GetFeature(ctx, featureName)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

func (s *OneStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.Feature, error) {
	features, err := s.db.ListFeature(ctx, opt)
	if err != nil {
		return nil, err
	}
	return features, nil
}

func (s *OneStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	return s.db.UpdateFeature(ctx, opt)
}
