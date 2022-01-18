package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Create metadata of a feature group.
func (s *OomStore) CreateGroup(ctx context.Context, opt types.CreateGroupOpt) (int, error) {
	entity, err := s.metadata.GetEntityByName(ctx, opt.EntityName)
	if err != nil {
		return 0, err
	}

	id, err := s.metadata.CreateGroup(ctx, metadata.CreateGroupOpt{
		GroupName:        opt.GroupName,
		EntityID:         entity.ID,
		Description:      opt.Description,
		Category:         opt.Category,
		SnapshotInterval: opt.SnapshotInterval,
	})
	if err != nil {
		return 0, err
	}

	if opt.Category == types.CategoryStream {
		if err = s.online.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
			Entity:  entity,
			GroupID: id,
		}); err != nil {
			return 0, err
		}
	}
	return id, nil
}

// Get metadata of a feature group by ID.
func (s *OomStore) GetGroup(ctx context.Context, id int) (*types.Group, error) {
	return s.metadata.GetGroup(ctx, id)
}

// Get metadata of a feature group by name.
func (s *OomStore) GetGroupByName(ctx context.Context, name string) (*types.Group, error) {
	return s.metadata.GetGroupByName(ctx, name)
}

// List metadata of feature groups under the same entity.
func (s *OomStore) ListGroup(ctx context.Context, entityID *int) (types.GroupList, error) {
	return s.metadata.ListGroup(ctx, entityID, nil)
}

// Update metadata of a feature group.
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
