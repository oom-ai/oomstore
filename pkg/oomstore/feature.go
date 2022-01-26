package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/online/sqlutil"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/util"
)

// GetFeature gets metadata of a feature by ID.
func (s *OomStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return s.metadata.GetFeature(ctx, id)
}

// GetFeatureByFullName gets metadata of a feature by full name.
func (s *OomStore) GetFeatureByFullName(ctx context.Context, fullName string) (*types.Feature, error) {
	groupName, featureName, err := util.SplitFullFeatureName(fullName)
	if err != nil {
		return nil, err
	}
	return s.GetFeatureByName(ctx, groupName, featureName)
}

// GetFeatureByName gets metadata of a feature by group name and feature name.
func (s *OomStore) GetFeatureByName(ctx context.Context, groupName string, featureName string) (*types.Feature, error) {
	return s.metadata.GetFeatureByName(ctx, groupName, featureName)
}

// ListFeature lists metadata of features meeting particular criteria.
func (s *OomStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	metadataOpt := metadata.ListFeatureOpt{}
	if opt.EntityName != nil {
		entity, err := s.metadata.GetEntityByName(ctx, *opt.EntityName)
		if err != nil {
			return nil, err
		}
		metadataOpt.EntityID = &entity.ID
	}
	if opt.GroupName != nil {
		group, err := s.metadata.GetGroupByName(ctx, *opt.GroupName)
		if err != nil {
			return nil, err
		}
		metadataOpt.GroupID = &group.ID
	}
	features, err := s.metadata.ListFeature(ctx, metadataOpt)
	if err != nil {
		return nil, err
	}
	if opt.FeatureNames != nil {
		features = features.FilterFullnames(*opt.FeatureNames)
	}
	return features, nil
}

// UpdateFeature updates metadata of a feature.
func (s *OomStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	if err := util.ValidateFullFeatureNames(opt.FeatureName); err != nil {
		return err
	}
	feature, err := s.GetFeatureByFullName(ctx, opt.FeatureName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
		FeatureID:      feature.ID,
		NewDescription: opt.NewDescription,
	})
}

// CreateFeature creates metadata of a feature.
func (s *OomStore) CreateFeature(ctx context.Context, opt types.CreateFeatureOpt) (int, error) {
	group, err := s.metadata.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return 0, err
	}

	revisions, err := s.metadata.ListRevision(ctx, &group.ID)
	if err != nil {
		return 0, err
	}
	if len(revisions) > 0 {
		return 0, errdefs.Errorf("group %s already has data and cannot add features due to the join and export mechanism", group.Name)
	}

	id, err := s.metadata.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: opt.FeatureName,
		GroupID:     group.ID,
		ValueType:   opt.ValueType,
		Description: opt.Description,
	})
	if err != nil {
		return 0, err
	}

	if group.Category == types.CategoryStream {
		features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
			GroupID: &group.ID,
		})
		if err != nil {
			return 0, err
		}
		if err = s.online.CreateTable(ctx, online.CreateTableOpt{
			EntityName: group.Entity.Name,
			TableName:  sqlutil.OnlineStreamTableName(group.ID),
			Features:   features,
		}); err != nil {
			return 0, err
		}
	}

	return id, nil
}
