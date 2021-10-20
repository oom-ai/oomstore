package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) ListRevision(ctx context.Context, groupName *string) ([]*types.Revision, error) {
	return s.db.ListRevision(ctx, groupName)
}
