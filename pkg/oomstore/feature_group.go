package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) CreateFeatureGroup(ctx context.Context, opt types.CreateFeatureGroupOpt) (int, error) {
	entity, err := s.metadata.GetEntityByName(ctx, opt.EntityName)
	if err != nil {
		return 0, err
	}
	return s.metadata.CreateFeatureGroup(ctx, metadata.CreateFeatureGroupOpt{
		GroupName:   opt.GroupName,
		EntityID:    entity.ID,
		Description: opt.Description,
		// Via the oomstore API, we can only create a batch feature group
		// So we hardcode the category to be batch
		Category: types.BatchFeatureCategory,
	})
}

func (s *OomStore) GetFeatureGroup(ctx context.Context, id int) (*types.FeatureGroup, error) {
	return s.metadata.GetFeatureGroup(ctx, id)
}

func (s *OomStore) GetFeatureGroupByName(ctx context.Context, name string) (*types.FeatureGroup, error) {
	return s.metadata.GetFeatureGroupByName(ctx, name)
}

func (s *OomStore) ListFeatureGroup(ctx context.Context, entityID *int) types.FeatureGroupList {
	return s.metadata.ListFeatureGroup(ctx, entityID)
}

func (s *OomStore) UpdateFeatureGroup(ctx context.Context, opt types.UpdateFeatureGroupOpt) error {
	group, err := s.metadata.GetFeatureGroupByName(ctx, opt.GroupName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateFeatureGroup(ctx, metadata.UpdateFeatureGroupOpt{
		GroupID:             group.ID,
		NewDescription:      opt.NewDescription,
		NewOnlineRevisionID: opt.NewOnlineRevisionID,
	})
}
