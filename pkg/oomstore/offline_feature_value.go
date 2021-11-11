package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OomStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) (*types.JoinResult, error) {
	features := s.metadatav2.ListFeature(ctx, metadatav2.ListFeatureOpt{FeatureIDs: opt.FeatureIDs})

	features = features.Filter(func(f *typesv2.Feature) bool {
		return f.Group.Category == types.BatchFeatureCategory
	})
	if len(features) == 0 {
		return nil, nil
	}

	entity, err := s.getSharedEntity(ctx, features)
	if err != nil || entity == nil {
		return nil, err
	}

	featureMap := buildGroupToFeaturesMap(features)
	revisionRangeMap := make(map[string][]*metadatav2.RevisionRange)
	for groupName := range featureMap {
		// TODO: This is slow but I haven't figured out a better way
		group, err := s.metadatav2.GetFeatureGroupByName(ctx, groupName)
		if err != nil {
			return nil, err
		}
		revisionRanges := s.metadatav2.BuildRevisionRanges(ctx, group.ID)
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
func buildGroupToFeaturesMap(features typesv2.FeatureList) map[string]typesv2.FeatureList {
	groups := make(map[string]typesv2.FeatureList)
	for _, f := range features {
		if _, ok := groups[f.Group.Name]; !ok {
			groups[f.Group.Name] = typesv2.FeatureList{}
		}
		groups[f.Group.Name] = append(groups[f.Group.Name], f)
	}
	return groups
}
