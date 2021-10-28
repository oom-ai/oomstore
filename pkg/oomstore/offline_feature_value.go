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
	features, err := s.metadata.ListRichFeature(ctx, types.ListFeatureOpt{FeatureNames: opt.FeatureNames})
	if err != nil {
		return nil, err
	}
	batchFeatures := features.Filter(func(f *types.RichFeature) bool {
		return f.Category == types.BatchFeatureCategory
	})

	// group_name -> []features
	featureGroups := buildGroupToFeaturesMap(batchFeatures)

	entityDataMap := make(map[string]dbutil.RowMap)
	for _, richFeatures := range featureGroups {
		if len(richFeatures) == 0 {
			continue
		}
		revisionRanges, err := s.metadata.BuildRevisionRanges(ctx, richFeatures[0].GroupName)
		if err != nil {
			return nil, err
		}
		entity, err := s.metadata.GetEntity(ctx, richFeatures[0].EntityName)
		if err != nil {
			return nil, err
		}
		featureValues, err := s.offline.Join(ctx, offline.JoinOpt{
			Entity:         entity,
			EntityRows:     opt.EntityRows,
			RevisionRanges: revisionRanges,
			Features:       richFeatures,
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
	for _, entity := range opt.EntityRows {
		key := entity.EntityKey + "," + strconv.Itoa(int(entity.UnixTime))
		if _, ok := entityDataMap[key]; !ok {
			entityDataMap[key] = dbutil.RowMap{
				"entity_key": entity.EntityKey,
				"unix_time":  entity.UnixTime,
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
func buildGroupToFeaturesMap(features types.RichFeatureList) map[string]types.RichFeatureList {
	groups := make(map[string]types.RichFeatureList)
	for _, f := range features {
		if _, ok := groups[f.GroupName]; !ok {
			groups[f.GroupName] = types.RichFeatureList{}
		}
		groups[f.GroupName] = append(groups[f.GroupName], f)
	}
	return groups
}
