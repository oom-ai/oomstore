package oomstore

import (
	"context"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Push inserts stream feature values to online store and offline store
func (s *OomStore) Push(ctx context.Context, opt types.PushOpt) error {
	features := s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		GroupName: &opt.GroupName,
	})
	var featureNames []string
	for name := range opt.FeatureValues {
		featureNames = append(featureNames, name)
	}

	entity := features[0].Entity()
	group := features[0].Group
	if group.Category != types.CategoryStream {
		return errdefs.Errorf("Push API is for streaming features only")
	}
	if !stringSliceEqual(features.Names(), featureNames) {
		return errdefs.Errorf("FeatureNames %v does not match with group's features %v", featureNames, features.Names())
	}
	values := make([]interface{}, 0, len(features))
	for _, f := range features {
		values = append(values, opt.FeatureValues[f.Name])
	}

	if err := s.online.Push(ctx, online.PushOpt{
		EntityName:    entity.Name,
		EntityKey:     opt.EntityKey,
		GroupID:       group.ID,
		Features:      features,
		FeatureValues: values,
	}); err != nil {
		return err
	}

	s.pushProcessor.Push(types.StreamRecord{
		GroupID:   group.ID,
		EntityKey: opt.EntityKey,
		UnixMilli: time.Now().UnixMilli(),
		Values:    values,
	})
	return nil
}
