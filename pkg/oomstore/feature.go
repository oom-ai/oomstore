package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) GetFeature(ctx context.Context, id int16) (*typesv2.Feature, error) {
	return s.metadatav2.GetFeature(ctx, id)
}

func (s *OomStore) GetFeatureByName(ctx context.Context, name string) (*typesv2.Feature, error) {
	return s.metadatav2.GetFeatureByName(ctx, name)
}

func (s *OomStore) ListFeature(ctx context.Context, opt metadatav2.ListFeatureOpt) typesv2.FeatureList {
	return s.metadatav2.ListFeature(ctx, opt)
}

func (s *OomStore) UpdateFeature(ctx context.Context, opt metadatav2.UpdateFeatureOpt) error {
	return s.metadatav2.UpdateFeature(ctx, opt)
}

func (s *OomStore) CreateBatchFeature(ctx context.Context, opt metadatav2.CreateFeatureOpt) (int16, error) {
	group, err := s.metadatav2.GetFeatureGroup(ctx, opt.GroupID)
	if err != nil {
		return 0, err
	}
	if group.Category != types.BatchFeatureCategory {
		return 0, fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}
	return s.metadatav2.CreateFeature(ctx, opt)
}
