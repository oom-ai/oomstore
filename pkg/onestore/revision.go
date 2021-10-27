package onestore

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/onestore/types"
)

func (s *OneStore) ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error) {
	return s.metadata.ListRevision(ctx, groupName)
}

func (s *OneStore) GetRevision(ctx context.Context, groupName string, revision int64) (*types.Revision, error) {
	return s.metadata.GetRevision(ctx, groupName, revision)
}
