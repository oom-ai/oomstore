package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

// GetFeature: get feature by featureName
func (s *OneStore) GetFeature(ctx context.Context, featureName string) (*types.Feature, error) {
	feature, err := s.db.GetFeature(ctx, featureName)
	if err != nil {
		return nil, err
	}
	return feature, nil
}
