package oomstore

import (
	"context"
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Push inserts stream feature values to online store and offline store
func (s *OomStore) Push(ctx context.Context, opt types.PushOpt) error {
	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		GroupName: &opt.GroupName,
	})
	if err != nil {
		return err
	}
	entity := features[0].Entity()
	group := features[0].Group
	if !stringSliceEqual(features.Names(), opt.FeatureNames) {
		return errdefs.Errorf("FeatureNames %v does not match with group's features %v", opt.FeatureNames, features.Names())
	}
	values := make([]interface{}, 0, len(features))
	for _, f := range features {
		for i, name := range opt.FeatureNames {
			if f.Name == name {
				values = append(values, opt.FeatureValues[i])
				break
			}
		}
	}

	if err = s.online.Push(ctx, online.PushOpt{
		Entity:        entity,
		EntityKey:     opt.EntityKey,
		GroupID:       group.ID,
		Features:      features,
		FeatureValues: values,
	}); err != nil {
		return err
	}

	s.streamPushProcessor.Push(types.StreamRecord{
		GroupID:   group.ID,
		EntityKey: opt.EntityKey,
		UnixMilli: time.Now().UnixMilli(),
		Values:    values,
	})
	return nil
}
