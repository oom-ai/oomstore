package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OomStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) (*types.JoinResult, error) {
	features, err := s.metadata.ListFeature(ctx, types.ListFeatureOpt{FeatureNames: opt.FeatureNames})
	if err != nil {
		return nil, err
	}
	features = features.Filter(func(f *types.Feature) bool {
		return f.Category == types.BatchFeatureCategory
	})
	if len(features) == 0 {
		return nil, nil
	}

	entityName, err := getEntityName(features)
	if err != nil || entityName == nil {
		return nil, err
	}
	entity, err := s.metadata.GetEntity(ctx, *entityName)
	if err != nil {
		return nil, err
	}

	featureMap := buildGroupToFeaturesMap(features)
	revisionRangeMap := make(map[string][]*types.RevisionRange)
	for groupName := range featureMap {
		revisionRanges, err := s.metadata.BuildRevisionRanges(ctx, groupName)
		if err != nil {
			return nil, err
		}
		revisionRangeMap[groupName] = revisionRanges
	}

	return s.offline.Join(ctx, offline.JoinOpt{
		Entity:           *entity,
		EntityRows:       opt.EntityRows,
		FeatureMap:       featureMap,
		RevisionRangeMap: revisionRangeMap,
	})
}

// key: group_name, value: slice of features
func buildGroupToFeaturesMap(features types.FeatureList) map[string]types.FeatureList {
	groups := make(map[string]types.FeatureList)
	for _, f := range features {
		if _, ok := groups[f.GroupName]; !ok {
			groups[f.GroupName] = types.FeatureList{}
		}
		groups[f.GroupName] = append(groups[f.GroupName], f)
	}
	return groups
}
