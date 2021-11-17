package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) ListRevision(ctx context.Context, groupID *int) types.RevisionList {
	return s.metadata.ListRevision(ctx, metadata.ListRevisionOpt{GroupID: groupID})
}

func (s *OomStore) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	return s.metadata.GetRevision(ctx, id)
}

func (s *OomStore) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	return s.metadata.GetRevisionBy(ctx, groupID, revision)
}
