package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

// GetFeature: get richFeature by featureName
func (s *OneStore) GetRichFeature(ctx context.Context, featureName string) (*types.RichFeature, error) {
	feature, err := s.db.GetRichFeature(ctx, featureName)
	if err != nil {
		return nil, err
	}
	return feature, nil
}

func (s *OneStore) ListRichFeature(ctx context.Context, opt types.ListFeatureOpt) ([]*types.RichFeature, error) {
	features, err := s.db.ListRichFeature(ctx, opt)
	if err != nil {
		return nil, err
	}
	return features, nil
}
