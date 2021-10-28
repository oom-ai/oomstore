package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) GetRichFeatureGroup(ctx context.Context, groupName string) (*types.RichFeatureGroup, error) {
	return s.metadata.GetRichFeatureGroup(ctx, groupName)
}

func (s *OomStore) ListRichFeatureGroup(ctx context.Context, entityName *string) ([]*types.RichFeatureGroup, error) {
	return s.metadata.ListRichFeatureGroup(ctx, entityName)
}
