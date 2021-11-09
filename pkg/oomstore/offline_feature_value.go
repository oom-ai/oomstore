package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OomStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) (<-chan *types.EntityRowWithFeatures, error) {
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

	joined, err := s.offline.Join(ctx, offline.JoinOpt{
		Entity:           *entity,
		EntityRows:       opt.EntityRows,
		FeatureMap:       featureMap,
		RevisionRangeMap: revisionRangeMap,
	})
	if err != nil {
		return nil, err
	}
	if joined == nil {
		return nil, nil
	}

	stream := make(chan *types.EntityRowWithFeatures)
	var processErr error
	go func() {
		defer close(stream)
		for item := range joined {
			if item.Error != nil {
				processErr = item.Error
				return
			}
			entityRowWithFeatures, tmpErr := processRowMap(item.RowMap, opt.FeatureNames)
			if tmpErr != nil {
				processErr = tmpErr
				return
			}
			stream <- entityRowWithFeatures
		}
	}()
	if processErr != nil {
		return nil, processErr
	}

	return stream, nil
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

func processRowMap(rowMap dbutil.RowMap, featureNames []string) (*types.EntityRowWithFeatures, error) {
	entityKey := cast.ToString(rowMap["entity_key"])
	unixTime := rowMap["unix_time"]
	unixTimeInt, err := castToInt64(unixTime)
	if err != nil {
		return nil, err
	}
	featureValues := make([]types.FeatureKV, 0, len(rowMap))
	for _, fn := range featureNames {
		featureValues = append(featureValues, types.NewFeatureKV(fn, rowMap[fn]))
	}
	return &types.EntityRowWithFeatures{
		EntityRow: types.EntityRow{
			EntityKey: entityKey,
			UnixTime:  unixTimeInt,
		},
		FeatureValues: featureValues,
	}, nil
}
