package onestore

import (
	"context"
	"fmt"

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
	richFeatures, err := s.db.ListRichFeature(ctx, opt)
	if err != nil {
		return nil, err
	}
	features := make([]*types.Feature, 0, len(richFeatures))
	for _, rf := range richFeatures {
		features = append(features, rf.ToFeature())
	}
	return features, nil
}

func (s *OneStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	return s.db.UpdateFeature(ctx, opt)
}

func (s *OneStore) CreateBatchFeature(ctx context.Context, opt types.CreateFeatureOpt) (*types.Feature, error) {
	group, err := s.db.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return nil, err
	}
	if group.Category != types.BatchFeatureCategory {
		return nil, fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}
	if err := s.db.CreateFeature(ctx, opt); err != nil {
		return nil, err
	}
	return s.db.GetFeature(ctx, opt.FeatureName)
}
