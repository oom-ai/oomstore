package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/metadata"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/util"
)

// OnlineGet gets online features of a particular entity instance.
func (s *OomStore) OnlineGet(ctx context.Context, opt types.OnlineGetOpt) (*types.FeatureValues, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}
	if opt.FeatureNames != nil {
		return s.OnlineGetByFeatures(ctx, opt.FeatureNames, opt.EntityKey)
	} else {
		return s.OnlineGetByGroup(ctx, *opt.GroupName, opt.EntityKey)
	}
}

func (s *OomStore) OnlineGetByGroup(ctx context.Context, groupName string, entityKey string) (*types.FeatureValues, error) {
	group, err := s.metadata.GetCachedGroupByName(ctx, groupName)
	if err != nil {
		return nil, err
	}
	featureValues, err := s.online.GetByGroup(ctx, online.GetByGroupOpt{
		EntityKey: entityKey,
		Group:     *group,
		ListFeature: func(groupID int) types.FeatureList {
			return s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
				GroupID: &groupID,
			})
		},
		RevisionID: group.OnlineRevisionID,
	})
	if err != nil {
		return nil, err
	}

	featureNames := make([]string, len(featureValues))
	featureValueMap := make(map[string]interface{})
	for featureName, featureValue := range featureValues {
		featureNames = append(featureNames, featureName)
		featureValueMap[featureName] = featureValue
	}
	return &types.FeatureValues{
		EntityName:      group.Entity.Name,
		EntityKey:       entityKey,
		FeatureNames:    featureNames,
		FeatureValueMap: featureValueMap,
	}, nil
}

func (s *OomStore) OnlineGetByFeatures(ctx context.Context, featureNames []string, entityKey string) (*types.FeatureValues, error) {
	rs := types.FeatureValues{
		EntityKey:       entityKey,
		FeatureNames:    featureNames,
		FeatureValueMap: make(map[string]interface{}),
	}
	features := s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		FullNames: &featureNames,
	})
	if len(features) == 0 {
		return &rs, nil
	}

	entity, err := features.GetSharedEntity()
	if err != nil {
		return nil, err
	}
	rs.EntityName = entity.Name

	featureMap := features.GroupByGroupID()

	for _, features := range featureMap {
		if len(features) == 0 {
			continue
		}

		group := features[0].Group
		if group.Category == types.CategoryBatch && group.OnlineRevisionID == nil {
			return nil, errdefs.Errorf("group %s has nil,please make sure you have already sync", group.Name)
		}

		featureValues, err := s.online.Get(ctx, online.GetOpt{
			EntityKey:  entityKey,
			Group:      *group,
			Features:   features,
			RevisionID: group.OnlineRevisionID,
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
	if err := util.ValidateFullFeatureNames(opt.FeatureNames...); err != nil {
		return nil, err
	}
	result := make(map[string]*types.FeatureValues)
	features := s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		FullNames: &opt.FeatureNames,
	})
	if len(features) == 0 {
		return result, nil
	}

	entity, err := features.GetSharedEntity()
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errdefs.Errorf("failed to get shared entity")
	}
	featureMap := features.GroupByGroupID()

	// entity_key -> feature_name -> feature_value
	featureValueMap, err := s.getFeatureValueMap(ctx, opt.EntityKeys, featureMap, entity)
	if err != nil {
		return nil, err
	}

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

	for _, features := range featureMap {
		if len(features) == 0 {
			continue
		}

		group := features[0].Group
		if group.Category == types.CategoryBatch && group.OnlineRevisionID == nil {
			return nil, errdefs.Errorf("group %s has nil,please make sure you have already sync", group.Name)
		}

		featureValues, err := s.online.MultiGet(ctx, online.MultiGetOpt{
			EntityKeys: entityKeys,
			Group:      *group,
			Features:   features,
			RevisionID: group.OnlineRevisionID,
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
