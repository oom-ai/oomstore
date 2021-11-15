package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) CreateFeatureGroup(ctx context.Context, opt metadata.CreateFeatureGroupOpt) (int16, error) {
	// Via the oomstore API, we can only create a batch feature group
	// So we hardcode the category to be batch
	opt.Category = types.BatchFeatureCategory
	return s.metadata.CreateFeatureGroup(ctx, opt)
}

func (s *OomStore) GetFeatureGroup(ctx context.Context, id int16) (*typesv2.FeatureGroup, error) {
	return s.metadata.GetFeatureGroup(ctx, id)
}

func (s *OomStore) GetFeatureGroupByName(ctx context.Context, name string) (*typesv2.FeatureGroup, error) {
	return s.metadata.GetFeatureGroupByName(ctx, name)
}

func (s *OomStore) ListFeatureGroup(ctx context.Context, entityID *int16) typesv2.FeatureGroupList {
	return s.metadata.ListFeatureGroup(ctx, entityID)
}

func (s *OomStore) UpdateFeatureGroup(ctx context.Context, opt metadata.UpdateFeatureGroupOpt) error {
	return s.metadata.UpdateFeatureGroup(ctx, opt)
}
