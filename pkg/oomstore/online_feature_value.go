package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) GetOnlineFeatureValues(ctx context.Context, opt types.GetOnlineFeatureValuesOpt) (types.FeatureValueMap, error) {
	m := make(map[string]interface{})

	features := s.metadatav2.ListFeature(ctx, metadata.ListFeatureOpt{FeatureNames: &opt.FeatureNames})
	features = features.Filter(func(f *typesv2.Feature) bool {
		return f.Group.OnlineRevisionID != nil
	})
	if len(features) == 0 {
		return m, nil
	}

	entity, err := s.getSharedEntity(features)
	if err != nil {
		return m, err
	}
	if entity == nil {
		return m, fmt.Errorf("failed to get shared entity")
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
			return m, err
		}
		for featureName, featureValue := range featureValues {
			m[featureName] = featureValue
		}
	}
	m[entity.Name] = opt.EntityKey
	return m, nil
}

func (s *OomStore) MultiGetOnlineFeatureValues(ctx context.Context, opt types.MultiGetOnlineFeatureValuesOpt) (types.FeatureDataSet, error) {
	features := s.metadatav2.ListFeature(ctx, metadata.ListFeatureOpt{FeatureIDs: &opt.FeatureIDs})

	features = features.Filter(func(f *typesv2.Feature) bool {
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

func (s *OomStore) getFeatureValueMap(ctx context.Context, entityKeys []string, featureMap map[int32]typesv2.FeatureList, entity *typesv2.Entity) (map[string]dbutil.RowMap, error) {
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
			featureValueMap[entityKey][entity.Name] = entityKey
		}
	}
	return featureValueMap, nil
}

func groupFeaturesByRevisionId(features typesv2.FeatureList) map[int32]typesv2.FeatureList {
	featureMap := make(map[int32]typesv2.FeatureList)
	for _, f := range features {
		id := f.OnlineRevisionID()
		if id == nil {
			continue
		}
		featureMap[*id] = append(featureMap[*id], f)
	}
	return featureMap
}

func (s *OomStore) getSharedEntity(features typesv2.FeatureList) (*typesv2.Entity, error) {
	m := make(map[int16]*typesv2.Entity)
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
