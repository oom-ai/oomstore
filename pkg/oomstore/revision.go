package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// List metadata of revisions of a same group.
func (s *OomStore) ListRevision(ctx context.Context, groupID *int) types.RevisionList {
	_ = s.metadata.Refresh()
	return s.metadata.CacheListRevision(ctx, groupID)
}

// Get metadata of a revision by ID.
func (s *OomStore) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CacheGetRevision(ctx, id)
}

// Get metadata of a revision by group ID and revision.
func (s *OomStore) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CacheGetRevisionBy(ctx, groupID, revision)
}
