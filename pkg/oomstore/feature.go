package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return s.metadata.GetFeature(ctx, id)
}

func (s *OomStore) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	return s.metadata.GetFeatureByName(ctx, name)
}

func (s *OomStore) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) types.FeatureList {
	return s.metadata.ListFeature(ctx, opt)
}

func (s *OomStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	feature, err := s.metadata.GetFeatureByName(ctx, opt.FeatureName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
		FeatureID:      feature.ID,
		NewDescription: opt.NewDescription,
	})
}

func (s *OomStore) CreateBatchFeature(ctx context.Context, opt types.CreateFeatureOpt) (int, error) {
	group, err := s.metadata.GetFeatureGroup(ctx, opt.GroupID)
	if err != nil {
		return 0, err
	}
	if group.Category != types.BatchFeatureCategory {
		return 0, fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}

	valueType, err := s.offline.TypeTag(opt.DBValueType)
	if err != nil {
		return 0, err
	}
	return s.metadata.CreateFeature(ctx, metadata.CreateFeatureOpt{
		CreateFeatureOpt: opt,
		ValueType:        valueType,
	})
}
