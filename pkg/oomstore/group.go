package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) CreateGroup(ctx context.Context, opt types.CreateGroupOpt) (int, error) {
	entity, err := s.metadata.GetEntityByName(ctx, opt.EntityName)
	if err != nil {
		return 0, err
	}
	return s.metadata.CreateGroup(ctx, metadata.CreateGroupOpt{
		GroupName:   opt.GroupName,
		EntityID:    entity.ID,
		Description: opt.Description,
		// Via the oomstore API, we can only create a batch feature group
		// So we hardcode the category to be batch
		Category: types.BatchFeatureCategory,
	})
}

func (s *OomStore) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	return s.metadata.GetGroup(ctx, id)
}

func (s *OomStore) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	return s.metadata.GetGroupByName(ctx, name)
}

func (s *OomStore) ListGroup(ctx context.Context, entityID *int) types.GroupList {
	return s.metadata.ListGroup(ctx, entityID)
}

func (s *OomStore) UpdateGroup(ctx context.Context, opt types.UpdateGroupOpt) error {
	group, err := s.metadata.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateGroup(ctx, metadata.UpdateGroupOpt{
		GroupID:             group.ID,
		NewDescription:      opt.NewDescription,
		NewOnlineRevisionID: opt.NewOnlineRevisionID,
	})
}
