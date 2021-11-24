package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Get online features of a particular entity instance.
func (s *OomStore) OnlineGet(ctx context.Context, opt types.OnlineGetOpt) (*types.FeatureValues, error) {
	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureNames: &opt.FeatureNames,
	})
	entity, err := getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	features = features.Filter(func(f *types.Feature) bool {
		return f.Group.OnlineRevisionID != nil
	})

	rs := types.FeatureValues{
		EntityName:      entity.Name,
		EntityKey:       opt.EntityKey,
		FeatureNames:    opt.FeatureNames,
		FeatureValueMap: make(map[string]interface{}),
	}

	if len(features) == 0 {
		return &rs, nil
	}

	featureMap := groupFeaturesByRevisionID(features)

	for onlineRevisionID, features := range featureMap {
		if len(features) == 0 {
			continue
		}
		featureValues, err := s.online.Get(ctx, online.GetOpt{
			Entity:      entity,
			RevisionID:  onlineRevisionID,
			EntityKey:   opt.EntityKey,
			FeatureList: features,
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

// Get online features of multiple entity instances.
func (s *OomStore) OnlineMultiGet(ctx context.Context, opt types.OnlineMultiGetOpt) (map[string]*types.FeatureValues, error) {
	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureNames: &opt.FeatureNames,
	})

	features = features.Filter(func(f *types.Feature) bool {
		return f.OnlineRevisionID() != nil
	})
	if len(features) == 0 {
		return nil, nil
	}

	entity, err := getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("failed to get shared entity")
	}
	featureMap := groupFeaturesByRevisionID(features)

	// entity_key -> feature_name -> feature_value
	featureValueMap, err := s.getFeatureValueMap(ctx, opt.EntityKeys, featureMap, entity)
	if err != nil {
		return nil, err
	}

	result := make(map[string]*types.FeatureValues)
	for _, entityKey := range opt.EntityKeys {
		result[entityKey] = &types.FeatureValues{
			EntityName:      entity.Name,
			EntityKey:       entityKey,
			FeatureNames:    opt.FeatureNames,
			FeatureValueMap: make(map[string]interface{}),
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

	for onlineRevisionID, features := range featureMap {
		if len(features) == 0 {
			continue
		}
		featureValues, err := s.online.MultiGet(ctx, online.MultiGetOpt{
			Entity:      entity,
			RevisionID:  onlineRevisionID,
			EntityKeys:  entityKeys,
			FeatureList: features,
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

func groupFeaturesByRevisionID(features types.FeatureList) map[int]types.FeatureList {
	featureMap := make(map[int]types.FeatureList)
	for _, f := range features {
		id := f.OnlineRevisionID()
		if id == nil {
			continue
		}
		featureMap[*id] = append(featureMap[*id], f)
	}
	return featureMap
}

func getSharedEntity(features types.FeatureList) (*types.Entity, error) {
	m := make(map[int]*types.Entity)
	for _, f := range features {
		m[f.Group.EntityID] = f.Group.Entity
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("expected 1 entity, got %d entities", len(m))
	}

	for _, entity := range m {
		return entity, nil
	}
	return nil, fmt.Errorf("expected 1 entity, got 0")
}
