package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) GetOnlineFeatureValues(ctx context.Context, opt types.GetOnlineFeatureValuesOpt) (types.FeatureValueMap, error) {
	m := make(map[string]interface{})

	features, err := s.metadata.GetRichFeatures(ctx, opt.FeatureNames)
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
	revisionIds, err := s.getRevisionIds(ctx, dataTables)
	if err != nil {
		return m, err
	}

	for dataTable, features := range dataTables {
		if len(features) == 0 {
			continue
		}
		revisionId, ok := revisionIds[dataTable]
		if !ok {
			continue
		}
		featureValues, err := s.online.Get(ctx, online.GetOpt{
			DataTable:   dataTable,
			EntityName:  *entityName,
			RevisionId:  revisionId,
			EntityKey:   opt.EntityKey,
			FeatureList: features.ToFeatureList(),
		})
		if err != nil {
			return m, err
		}
		for featureName, featureValue := range featureValues {
			m[featureName] = featureValue
		}
	}
	return m, nil
}

func (s *OomStore) MultiGetOnlineFeatureValues(ctx context.Context, opt types.MultiGetOnlineFeatureValuesOpt) (*types.FeatureDataSet, error) {
	features, err := s.metadata.GetRichFeatures(ctx, opt.FeatureNames)
	if err != nil {
		return nil, err
	}
	features = filterAvailableFeatures(features)
	if len(features) == 0 {
		return nil, nil
	}

	entityName, err := getEntityName(features)
	if err != nil || entityName == nil {
		return nil, err
	}
	dataTables := getDataTables(features)
	revisionIds, err := s.getRevisionIds(ctx, dataTables)
	if err != nil {
		return nil, err
	}

	// entity_key -> feature_name -> feature_value
	featureValueMap, err := s.getFeatureValueMap(ctx, opt.EntityKeys, dataTables, revisionIds, *entityName)
	if err != nil {
		return nil, err
	}

	return buildFeatureDataSet(featureValueMap, opt)
}

func (s *OomStore) getFeatureValueMap(ctx context.Context, entityKeys []string, dataTableMap map[string]types.RichFeatureList, revisionIds map[string]int32, entityName string) (map[string]dbutil.RowMap, error) {
	// entity_key -> types.RecordMap
	featureValueMap := make(map[string]dbutil.RowMap)

	for dataTable, features := range dataTableMap {
		if len(features) == 0 {
			continue
		}
		revisionId, ok := revisionIds[dataTable]
		if !ok {
			continue
		}

		featureValues, err := s.online.MultiGet(ctx, online.MultiGetOpt{
			DataTable:   dataTable,
			EntityName:  entityName,
			RevisionId:  revisionId,
			EntityKeys:  entityKeys,
			FeatureList: features.ToFeatureList(),
		})
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

func (s *OomStore) getRevisionIds(ctx context.Context, dataTables map[string]types.RichFeatureList) (map[string]int32, error) {
	dataTableSlice := make([]string, 0, len(dataTables))
	for dataTable := range dataTables {
		dataTableSlice = append(dataTableSlice, dataTable)
	}
	revisions, err := s.metadata.GetRevisionsByDataTables(ctx, dataTableSlice)
	if err != nil {
		return nil, nil
	}
	revisionMap := make(map[string]int32)
	for _, revision := range revisions {
		revisionMap[revision.DataTable] = revision.ID
	}
	return revisionMap, nil
}
func filterAvailableFeatures(features types.RichFeatureList) (rs types.RichFeatureList) {
	for _, f := range features {
		if f.DataTable != nil {
			rs = append(rs, f)
		}
	}
	return
}

func getDataTables(features types.RichFeatureList) map[string]types.RichFeatureList {
	dataTableMap := make(map[string]types.RichFeatureList)
	for _, f := range features {
		dataTable := *f.DataTable
		dataTableMap[dataTable] = append(dataTableMap[dataTable], f)
	}
	return dataTableMap
}

func getEntityName(features types.RichFeatureList) (*string, error) {
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

func buildFeatureDataSet(valueMap map[string]dbutil.RowMap, opt types.MultiGetOnlineFeatureValuesOpt) (*types.FeatureDataSet, error) {
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
