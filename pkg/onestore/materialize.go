package onestore

import (
	"context"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) Materialize(ctx context.Context, opt types.MaterializeOpt) error {
	features, err := s.ListFeature(ctx, types.ListFeatureOpt{GroupName: &opt.GroupName})
	if err != nil {
		return err
	}

	featureNames := []string{}
	for _, f := range features {
		featureNames = append(featureNames, f.Name)
	}

	stream, err := s.offline.GetFeatureValuesStream(ctx, types.GetFeatureValuesStreamOpt{
		GroupName:    opt.GroupName,
		Revision:     opt.GroupRevision,
		FeatureNames: featureNames,
	})
	if err != nil {
		return err
	}

	revision, err := s.GetRevision(ctx, opt.GroupName, opt.GroupRevision)
	if err != nil {
		return err
	}

	if err := s.online.SinkFeatureValuesStream(ctx, stream, features, revision); err != nil {
		return err
	}

	return nil
}
