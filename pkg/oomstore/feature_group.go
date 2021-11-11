package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) CreateFeatureGroup(ctx context.Context, opt metadatav2.CreateFeatureGroupOpt) (int16, error) {
	return s.metadatav2.CreateFeatureGroup(ctx, opt)
}

func (s *OomStore) GetFeatureGroup(ctx context.Context, id int16) (*typesv2.FeatureGroup, error) {
	return s.metadatav2.GetFeatureGroup(ctx, id)
}

func (s *OomStore) GetFeatureGroupByName(ctx context.Context, name string) (*typesv2.FeatureGroup, error) {
	return s.metadatav2.GetFeatureGroupByName(ctx, name)
}

func (s *OomStore) ListFeatureGroup(ctx context.Context, entityID *int16) typesv2.FeatureGroupList {
	return s.metadatav2.ListFeatureGroup(ctx, entityID)
}

func (s *OomStore) UpdateFeatureGroup(ctx context.Context, opt metadatav2.UpdateFeatureGroupOpt) error {
	return s.metadatav2.UpdateFeatureGroup(ctx, opt)
}
