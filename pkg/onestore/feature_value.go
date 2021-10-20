package onestore

import (
	"context"
	"fmt"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) GetOnlineFeatureValues(ctx context.Context, opt types.GetOnlineFeatureValuesOpt) (types.FeatureValueMap, error) {
	m := make(map[string]interface{})

	features, err := s.db.GetRichFeatures(ctx, opt.FeatureNames)
	if err != nil {
		return m, err
	}
	features = filterAvailableFeatures(features)
	if len(features) == 0 {
		return m, nil
	}

	entityName, err := getEntityName(features)
	if err != nil || entityName == nil {
		return m, err
	}
	dataTables := getDataTables(features)

	for dataTable, featureNames := range dataTables {
		featureValues, err := s.db.GetFeatureValues(ctx, dataTable, *entityName, opt.EntityValue, featureNames)
		if err != nil {
			return m, err
		}
		for featureName, featureValue := range featureValues {
			m[featureName] = featureValue
		}
	}
	return m, nil
}

func filterAvailableFeatures(features []*types.RichFeature) (rs []*types.RichFeature) {
	for _, f := range features {
		if f.DataTable != nil {
			rs = append(rs, f)
		}
	}
	return
}

func getDataTables(features []*types.RichFeature) map[string][]string {
	dataTableMap := make(map[string][]string)
	for _, f := range features {
		dataTable := *f.DataTable
		dataTableMap[dataTable] = append(dataTableMap[dataTable], f.Name)
	}
	return dataTableMap
}

func getEntityName(features []*types.RichFeature) (*string, error) {
	m := make(map[string]string)
	for _, f := range features {
		m[f.EntityName] = f.Name
	}
	if len(m) > 1 {
		return nil, fmt.Errorf("inconsistent entity type: %v", m)
	}
	for entityName := range m {
		return &entityName, nil
	}
	return nil, nil
}

func (s *OneStore) GetOnlineFeatureValuesWithMultiEntityValues(ctx context.Context, opt types.GetOnlineFeatureValuesWithMultiEntityValuesOpt) (*types.FeatureDataSet, error) {
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
		featureValues, err := s.db.GetFeatureValuesWithMultiEntityValues(ctx, dataTable, entityName, entityKeys, featureNames)
		if err != nil {
			return nil, err
		}
		for entityKey, m := range featureValues {
			if featureValueMap[entityKey] == nil {
				featureValueMap[entityKey] = make(map[string]interface{})
			}
			for fn, fv := range m {
				featureValueMap[entityKey][fn] = fv
			}
		}
	}
	return featureValueMap, nil
}

func buildFeatureDataSet(valueMap map[string]database.RowMap, opt types.GetOnlineFeatureValuesWithMultiEntityValuesOpt) (*types.FeatureDataSet, error) {
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
			for fn, fv := range m {
				entityDataMap[key][fn] = fv
			}
		}
	}

	entityDataSet := make([]*types.EntityRowWithFeatures, 0, len(entityDataMap))
	for _, rowMap := range entityDataMap {
		entityKey := rowMap["entity_key"]
		unixTime := rowMap["unix_time"]
		delete(rowMap, "entity_key")
		delete(rowMap, "unix_time")

		featureValues := make([]types.FeatureKV, 0, len(rowMap))
		for fn, fv := range rowMap {
			featureValues = append(featureValues, types.NewFeatureKV(fn, fv))
		}
		entityDataSet = append(entityDataSet, &types.EntityRowWithFeatures{
			EntityRow: types.EntityRow{
				EntityKey: entityKey.(string),
				UnixTime:  unixTime.(int64),
			},
			FeatureValues: featureValues,
		})
	}

	return entityDataSet, nil
}
