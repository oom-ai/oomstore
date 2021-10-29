package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	return s.metadata.GetFeature(ctx, featureName)
}

func (s *OomStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	return s.metadata.ListFeature(ctx, opt)
}

func (s *OomStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	return s.metadata.UpdateFeature(ctx, opt)
}

func (s *OomStore) CreateBatchFeature(ctx context.Context, opt types.CreateFeatureOpt) error {
	valueType, err := s.offline.TypeTag(opt.DBValueType)
	if err != nil {
		return err
	}
	group, err := s.metadata.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return err
	}
	if group.Category != types.BatchFeatureCategory {
		return fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}
	return s.metadata.CreateFeature(ctx, metadata.CreateFeatureOpt{
		CreateFeatureOpt: opt,
		ValueType:        valueType,
	})
}
