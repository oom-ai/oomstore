package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error) {
	return s.metadata.ListRevision(ctx, metadata.ListRevisionOpt{GroupName: groupName})
}

func (s *OomStore) GetRevision(ctx context.Context, groupName string, revision int64) (*types.Revision, error) {
	return s.metadata.GetRevision(ctx, metadata.GetRevisionOpt{
		GroupName: &groupName,
		Revision:  &revision,
	})
}
