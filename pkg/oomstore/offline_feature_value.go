package oomstore

import (
	"context"
	"sort"
	"strconv"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OomStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) ([]*types.EntityRowWithFeatures, error) {
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

	// group_name -> []features
	featureGroups := buildGroupToFeaturesMap(features)

	entityDataMap := make(map[string]dbutil.RowMap)
	for groupName, featureSlice := range featureGroups {
		if len(featureSlice) == 0 {
			continue
		}
		revisionRanges, err := s.metadata.BuildRevisionRanges(ctx, groupName)
		if err != nil {
			return nil, err
		}
		featureValues, err := s.offline.Join(ctx, offline.JoinOpt{
			Entity:         entity,
			EntityRows:     opt.EntityRows,
			RevisionRanges: revisionRanges,
			Features:       featureSlice,
		})
		if err != nil {
			return nil, err
		}
		for key, m := range featureValues {
			if _, ok := entityDataMap[key]; !ok {
				entityDataMap[key] = make(dbutil.RowMap)
			}
			for fn, fv := range m {
				entityDataMap[key][fn] = fv
			}
		}
	}
	for _, e := range opt.EntityRows {
		key := e.EntityKey + "," + strconv.Itoa(int(e.UnixTime))
		if _, ok := entityDataMap[key]; !ok {
			entityDataMap[key] = dbutil.RowMap{
				"entity_key": e.EntityKey,
				"unix_time":  e.UnixTime,
			}
		}
	}

	entityDataSet := make([]*types.EntityRowWithFeatures, 0, len(entityDataMap))
	for _, rowMap := range entityDataMap {
		entityKey := rowMap["entity_key"]
		unixTime := rowMap["unix_time"]
		unixTimeInt, err := castToInt64(unixTime)
		if err != nil {
			return nil, err
		}

		featureValues := make([]types.FeatureKV, 0, len(rowMap))
		for _, fn := range opt.FeatureNames {
			featureValues = append(featureValues, types.NewFeatureKV(fn, rowMap[fn]))
		}
		entityDataSet = append(entityDataSet, &types.EntityRowWithFeatures{
			EntityRow: types.EntityRow{
				EntityKey: cast.ToString(entityKey),
				UnixTime:  unixTimeInt,
			},
			FeatureValues: featureValues,
		})
	}
	sort.Slice(entityDataSet, func(i, j int) bool {
		if entityDataSet[i].EntityKey == entityDataSet[j].EntityKey {
			return entityDataSet[i].UnixTime < entityDataSet[j].UnixTime
		}
		return entityDataSet[i].EntityKey < entityDataSet[j].EntityKey
	})
	return entityDataSet, nil
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
