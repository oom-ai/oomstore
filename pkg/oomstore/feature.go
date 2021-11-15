package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) GetFeature(ctx context.Context, id int16) (*typesv2.Feature, error) {
	return s.metadata.GetFeature(ctx, id)
}

func (s *OomStore) GetFeatureByName(ctx context.Context, name string) (*typesv2.Feature, error) {
	return s.metadata.GetFeatureByName(ctx, name)
}

func (s *OomStore) ListFeature(ctx context.Context, opt metadata.ListFeatureOpt) typesv2.FeatureList {
	return s.metadata.ListFeature(ctx, opt)
}

func (s *OomStore) UpdateFeature(ctx context.Context, opt metadata.UpdateFeatureOpt) error {
	return s.metadata.UpdateFeature(ctx, opt)
}

func (s *OomStore) CreateBatchFeature(ctx context.Context, opt metadata.CreateFeatureOpt) (int16, error) {
	group, err := s.metadata.GetFeatureGroup(ctx, opt.GroupID)
	if err != nil {
		return 0, err
	}
	if group.Category != types.BatchFeatureCategory {
		return 0, fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}

	if opt.ValueType, err = s.offline.TypeTag(opt.DBValueType); err != nil {
		return 0, err
	}
	return s.metadata.CreateFeature(ctx, opt)
}
