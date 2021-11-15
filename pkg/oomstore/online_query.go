package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) GetOnlineFeatureValues(ctx context.Context, opt types.GetOnlineFeatureValuesOpt) (*types.FeatureValues, error) {
	m := make(map[string]interface{})

	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{FeatureNames: &opt.FeatureNames})
	features = features.Filter(func(f *types.Feature) bool {
		return f.Group.OnlineRevisionID != nil
	})
	if len(features) == 0 {
		return nil, nil
	}

	entity, err := s.getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("failed to get shared entity")
	}
	featureMap := groupFeaturesByRevisionId(features)

	for onlineRevisionId, features := range featureMap {
		if len(features) == 0 {
			continue
		}
		featureValues, err := s.online.Get(ctx, online.GetOpt{
			Entity:      entity,
			RevisionID:  onlineRevisionId,
			EntityKey:   opt.EntityKey,
			FeatureList: features,
		})
		if err != nil {
			return nil, err
		}
		for featureName, featureValue := range featureValues {
			m[featureName] = featureValue
		}
	}
	return &types.FeatureValues{
		EntityName:      entity.Name,
		EntityKey:       opt.EntityKey,
		FeatureNames:    opt.FeatureNames,
		FeatureValueMap: m,
	}, nil
}

func (s *OomStore) MultiGetOnlineFeatureValues(ctx context.Context, opt types.MultiGetOnlineFeatureValuesOpt) (types.FeatureDataSet, error) {
	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{FeatureIDs: &opt.FeatureIDs})

	features = features.Filter(func(f *types.Feature) bool {
		return f.OnlineRevisionID() != nil
	})
	if len(features) == 0 {
		return nil, nil
	}

	entity, err := s.getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("failed to get shared entity")
	}
	featureMap := groupFeaturesByRevisionId(features)

	// entity_key -> feature_name -> feature_value
	featureValueMap, err := s.getFeatureValueMap(ctx, opt.EntityKeys, featureMap, entity)
	if err != nil {
		return nil, err
	}

	return buildFeatureDataSet(featureValueMap, features.Names(), opt.EntityKeys)
}

func (s *OomStore) getFeatureValueMap(ctx context.Context, entityKeys []string, featureMap map[int32]types.FeatureList, entity *types.Entity) (map[string]dbutil.RowMap, error) {
	// entity_key -> types.RecordMap
	featureValueMap := make(map[string]dbutil.RowMap)

	for onlineRevisionId, features := range featureMap {
		if len(features) == 0 {
			continue
		}
		featureValues, err := s.online.MultiGet(ctx, online.MultiGetOpt{
			Entity:      entity,
			RevisionID:  onlineRevisionId,
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

func groupFeaturesByRevisionId(features types.FeatureList) map[int32]types.FeatureList {
	featureMap := make(map[int32]types.FeatureList)
	for _, f := range features {
		id := f.OnlineRevisionID()
		if id == nil {
			continue
		}
		featureMap[*id] = append(featureMap[*id], f)
	}
	return featureMap
}

func (s *OomStore) getSharedEntity(features types.FeatureList) (*types.Entity, error) {
	m := make(map[int16]*types.Entity)
	for _, f := range features {
		m[f.Group.EntityID] = f.Group.Entity
	}
	if len(m) != 1 {
		return nil, fmt.Errorf("expected 1 entity, got %d entities", len(m))
	}

	for _, entity := range m {
		return entity, nil
	}
	return nil, nil
}

func buildFeatureDataSet(valueMap map[string]dbutil.RowMap, featureNames []string, entityKeys []string) (types.FeatureDataSet, error) {
	fds := types.NewFeatureDataSet()
	for _, entityKey := range entityKeys {
		fds[entityKey] = make([]types.FeatureKV, 0)
		// TODO: double check the logic doesn't change
		for _, fn := range featureNames {
			if fv, ok := valueMap[entityKey][fn]; ok {
				fds[entityKey] = append(fds[entityKey], types.NewFeatureKV(fn, fv))
			} else {
				return nil, fmt.Errorf("missing feature %s for entity %s", fn, entityKey)
			}
		}
	}
	return fds, nil
}
