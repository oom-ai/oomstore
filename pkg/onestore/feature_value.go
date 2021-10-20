package onestore

import (
	"context"
	"fmt"
	"sort"
	"strconv"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cast"
)

func (s *OneStore) GetOnlineFeatureValues(ctx context.Context, opt types.GetOnlineFeatureValuesOpt) (*types.FeatureDataSet, error) {
	features, err := s.db.GetRichFeatures(ctx, opt.FeatureNames)
	if err != nil {
		return nil, err
	}

	// data_table -> []feature_name
	dataTableMap := buildDataTableMap(features)
	// data_table -> entity_name
	entityNameMap := buildEntityNameMap(features)

	// entity_key -> feature_name -> feature_value
	featureValueMap, err := s.getFeatureValueMap(ctx, opt.EntityKeys, dataTableMap, entityNameMap)
	if err != nil {
		return nil, err
	}

	return buildFeatureDataSet(featureValueMap, opt)
}

func (s *OneStore) getFeatureValueMap(ctx context.Context, entityKeys []string, dataTableMap map[string][]string, entityNameMap map[string]string) (map[string]database.RowMap, error) {
	// entity_key -> types.RecordMap
	featureValueMap := make(map[string]database.RowMap)

	for dataTable, featureNames := range dataTableMap {
		entityName, ok := entityNameMap[dataTable]
		if !ok {
			return nil, fmt.Errorf("missing entity_name for table %s", dataTable)
		}
		featureValues, err := s.db.GetFeatureValues(ctx, dataTable, entityName, entityKeys, featureNames)
		if err != nil {
			return nil, err
		}
		for entityKey, m := range featureValues {
			for fn, fv := range m {
				featureValueMap[entityKey][fn] = fv
			}
		}
	}
	return featureValueMap, nil
}

func buildFeatureDataSet(valueMap map[string]database.RowMap, opt types.GetOnlineFeatureValuesOpt) (*types.FeatureDataSet, error) {
	fds := types.NewFeatureDataSet()
	for _, entityKey := range opt.EntityKeys {
		fds[entityKey] = make([]types.FeatureKV, 0)
		for _, fn := range opt.FeatureNames {
			if fv, ok := valueMap[entityKey][fn]; ok {
				fds[entityKey] = append(fds[entityKey], types.NewFeatureKV(fn, fv))
			} else {
				return nil, fmt.Errorf("missing feature %s for entity %s", fn, entityKey)
			}
		}
	}
	return &fds, nil
}

// key: data_table, value: slice of feature_names
func buildDataTableMap(features []*types.RichFeature) map[string][]string {
	dataTableMap := make(map[string][]string)
	for _, f := range features {
		if f.DataTable == nil {
			continue
		}
		dataTable := *f.DataTable
		if _, ok := dataTableMap[dataTable]; !ok {
			dataTableMap[dataTable] = make([]string, 0)
		}
		dataTableMap[dataTable] = append(dataTableMap[dataTable], f.Name)
	}
	return dataTableMap
}

// key: data_table, value: entity_name
func buildEntityNameMap(features []*types.RichFeature) map[string]string {
	entityNameMap := make(map[string]string)
	for _, f := range features {
		if f.DataTable == nil {
			continue
		}
		dataTable := *f.DataTable
		if _, ok := entityNameMap[dataTable]; !ok {
			entityNameMap[dataTable] = f.EntityName
		}
	}
	return entityNameMap
}

// key: data_table, value: slice of features
func buildDataTableToFeaturesMap(features []*types.RichFeature) map[string][]*types.RichFeature {
	dataTableToFeaturesMap := make(map[string][]*types.RichFeature)
	for _, f := range features {
		if f.DataTable == nil {
			continue
		}
		dataTable := *f.DataTable
		if _, ok := dataTableToFeaturesMap[dataTable]; !ok {
			dataTableToFeaturesMap[dataTable] = make([]*types.RichFeature, 0)
		}
		dataTableToFeaturesMap[dataTable] = append(dataTableToFeaturesMap[dataTable], f)
	}
	return dataTableToFeaturesMap
}

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OneStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) ([]*types.EntityRowWithFeatures, error) {
	features, err := s.db.GetRichFeatures(ctx, opt.FeatureNames)
	if err != nil {
		return nil, err
	}
	batchFeatures := make([]*types.RichFeature, 0)
	for _, f := range features {
		if f.Category == types.BatchFeatureCategory {
			batchFeatures = append(batchFeatures, f)
		}
	}
	// data_table -> []features
	dataTableToFeaturesMap := buildDataTableToFeaturesMap(batchFeatures)

	entityDataMap := make(map[string]database.RowMap)
	for _, richFeatures := range dataTableToFeaturesMap {
		featureValues, err := s.db.GetPointInTimeFeatureValues(ctx, richFeatures, opt.EntityRows)
		if err != nil {
			return nil, err
		}
		for key, m := range featureValues {
			if _, ok := entityDataMap[key]; !ok {
				entityDataMap[key] = make(database.RowMap)
			}
			for fn, fv := range m {
				entityDataMap[key][fn] = fv
			}
		}
	}
	for _, entity := range opt.EntityRows {
		key := entity.EntityKey + "," + strconv.Itoa(int(entity.UnixTime))
		if _, ok := entityDataMap[key]; !ok {
			entityDataMap[key] = database.RowMap{
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
