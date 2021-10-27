package onestore

import (
	"context"
	"sort"
	"strconv"

	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cast"
)

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OneStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) ([]*types.EntityRowWithFeatures, error) {
	features, err := s.metadata.GetRichFeatures(ctx, opt.FeatureNames)
	if err != nil {
		return nil, err
	}
	batchFeatures := make([]*types.RichFeature, 0)
	for _, f := range features {
		if f.Category == types.BatchFeatureCategory {
			batchFeatures = append(batchFeatures, f)
		}
	}
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
		featureValues, err := s.offline.GetPointInTimeFeatureValues(ctx, entity, opt.EntityRows, revisionRanges, richFeatures)
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
func buildGroupToFeaturesMap(features []*types.RichFeature) map[string][]*types.RichFeature {
	groups := make(map[string][]*types.RichFeature)
	for _, f := range features {
		if _, ok := groups[f.GroupName]; !ok {
			groups[f.GroupName] = make([]*types.RichFeature, 0)
		}
		groups[f.GroupName] = append(groups[f.GroupName], f)
	}
	return groups
}
