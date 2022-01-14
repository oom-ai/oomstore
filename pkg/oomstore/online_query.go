package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// OnlineGet gets online features of a particular entity instance.
func (s *OomStore) OnlineGet(ctx context.Context, opt types.OnlineGetOpt) (*types.FeatureValues, error) {
	if err := validateFullFeatureNames(opt.FeatureFullNames...); err != nil {
		return nil, err
	}
	rs := types.FeatureValues{
		EntityKey:        opt.EntityKey,
		FeatureFullNames: opt.FeatureFullNames,
		FeatureValueMap:  make(map[string]interface{}),
	}
	features := s.metadata.ListCachedFeature(ctx, &opt.FeatureFullNames)
	if len(features) == 0 {
		return &rs, nil
	}

	entity, err := getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	rs.EntityName = entity.Name

	featureMap := groupFeaturesByGroupID(features)

	for _, features := range featureMap {
		if len(features) == 0 {
			continue
		}
		revisionID := features[0].OnlineRevisionID()
		group := features[0].Group
		featureValues, err := s.online.Get(ctx, online.GetOpt{
			Entity:     entity,
			RevisionID: revisionID,
			Group:      group,
			EntityKey:  opt.EntityKey,
			Features:   features,
		})
		if err != nil {
			return nil, err
		}
		for featureName, featureValue := range featureValues {
			rs.FeatureValueMap[featureName] = featureValue
		}
	}
	return &rs, nil
}

// OnlineMultiGet gets online features of multiple entity instances.
func (s *OomStore) OnlineMultiGet(ctx context.Context, opt types.OnlineMultiGetOpt) (map[string]*types.FeatureValues, error) {
	if err := validateFullFeatureNames(opt.FeatureFullNames...); err != nil {
		return nil, err
	}
	result := make(map[string]*types.FeatureValues)
	features := s.metadata.ListCachedFeature(ctx, &opt.FeatureFullNames)
	if len(features) == 0 {
		return result, nil
	}

	entity, err := getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errdefs.Errorf("failed to get shared entity")
	}
	featureMap := groupFeaturesByGroupID(features)

	// entity_key -> feature_name -> feature_value
	featureValueMap, err := s.getFeatureValueMap(ctx, opt.EntityKeys, featureMap, entity)
	if err != nil {
		return nil, err
	}

	for _, entityKey := range opt.EntityKeys {
		result[entityKey] = &types.FeatureValues{
			EntityName:       entity.Name,
			EntityKey:        entityKey,
			FeatureFullNames: opt.FeatureFullNames,
			FeatureValueMap:  make(map[string]interface{}),
		}
		for featureName, featureValue := range featureValueMap[entityKey] {
			result[entityKey].FeatureValueMap[featureName] = featureValue
		}
	}
	return result, nil
}

func (s *OomStore) getFeatureValueMap(ctx context.Context, entityKeys []string, featureMap map[int]types.FeatureList, entity *types.Entity) (map[string]dbutil.RowMap, error) {
	// entity_key -> types.RecordMap
	featureValueMap := make(map[string]dbutil.RowMap)

	for _, features := range featureMap {
		if len(features) == 0 {
			continue
		}
		revisionID := features[0].OnlineRevisionID()
		group := features[0].Group
		featureValues, err := s.online.MultiGet(ctx, online.MultiGetOpt{
			Entity:     entity,
			RevisionID: revisionID,
			Group:      group,
			EntityKeys: entityKeys,
			Features:   features,
		})
		if err != nil {
			return nil, err
		}
		for entityKey, m := range featureValues {
			if featureValueMap[entityKey] == nil {
				featureValueMap[entityKey] = make(dbutil.RowMap)
			}
			for fn, fv := range m {
				featureValueMap[entityKey][fn] = fv
			}
		}
	}
	return featureValueMap, nil
}

func groupFeaturesByGroupID(features types.FeatureList) map[int]types.FeatureList {
	featureMap := make(map[int]types.FeatureList)
	for _, f := range features {
		id := f.GroupID
		featureMap[id] = append(featureMap[id], f)
	}
	return featureMap
}

func getSharedEntity(features types.FeatureList) (*types.Entity, error) {
	m := make(map[int]*types.Entity)
	for _, f := range features {
		m[f.Group.EntityID] = f.Group.Entity
	}
	if len(m) != 1 {
		return nil, errdefs.Errorf("expected 1 entity, got %d entities", len(m))
	}

	for _, entity := range m {
		return entity, nil
	}
	return nil, errdefs.Errorf("expected 1 entity, got 0")
}
