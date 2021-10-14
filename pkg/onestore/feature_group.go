package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) CreateGroup(ctx context.Context, opt types.CreateFeatureGroupOpt) error {
	return s.db.CreateFeatureGroup(ctx, opt)
}

func (s *OneStore) GetFeatureGroup(ctx context.Context, groupName string) (*types.FeatureGroup, error) {
	return s.db.GetFeatureGroup(ctx, groupName)
}
