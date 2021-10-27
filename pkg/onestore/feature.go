package onestore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

// GetFeature: get feature by featureName
func (s *OneStore) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	feature, err := s.metadata.GetFeature(ctx, featureName)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

func (s *OneStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.Feature, error) {
	richFeatures, err := s.metadata.ListRichFeature(ctx, opt)
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
	return s.metadata.UpdateFeature(ctx, opt)
}

func (s *OneStore) CreateBatchFeature(ctx context.Context, opt types.CreateFeatureOpt) (*types.Feature, error) {
	valueType, err := s.offline.TypeTag(opt.DBValueType)
	if err != nil {
		return nil, err
	}
	group, err := s.metadata.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return nil, err
	}
	if group.Category != types.BatchFeatureCategory {
		return nil, fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}
	if err := s.metadata.CreateFeature(ctx, metadata.CreateFeatureOpt{
		CreateFeatureOpt: opt,
		ValueType:        valueType,
	}); err != nil {
		return nil, err
	}
	return s.metadata.GetFeature(ctx, opt.FeatureName)
}
