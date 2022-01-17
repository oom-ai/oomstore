package oomstore

import (
	"context"
	"strings"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/util"
)

// Get metadata of a feature by ID.
func (s *OomStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return s.metadata.GetFeature(ctx, id)
}

// Get metadata of a feature by full name.
func (s *OomStore) GetFeatureByFullName(ctx context.Context, fullName string) (*types.Feature, error) {
	groupName, featureName, err := util.SplitFullFeatureName(fullName)
	if err != nil {
		return nil, err
	}
	return s.GetFeatureByName(ctx, groupName, featureName)
}

// Get metadata of a feature by group name and feature name.
func (s *OomStore) GetFeatureByName(ctx context.Context, groupName string, featureName string) (*types.Feature, error) {
	return s.metadata.GetFeatureByName(ctx, groupName, featureName)
}

// List metadata of features meeting particular criteria.
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

// Update metadata of a feature.
func (s *OomStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	if err := validateFullFeatureNames(opt.FeatureName); err != nil {
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

// Create metadata of a feature.
func (s *OomStore) CreateFeature(ctx context.Context, opt types.CreateFeatureOpt) (int, error) {
	group, err := s.metadata.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return 0, err
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
		feature, err := s.metadata.GetFeature(ctx, id)
		if err != nil {
			return 0, err
		}
		if err := s.online.PrepareStreamTable(ctx, online.PrepareStreamTableOpt{
			Entity:  group.Entity,
			GroupID: group.ID,
			Feature: feature,
		}); err != nil {
			return 0, err
		}
	}

	return id, nil
}

func validateFullFeatureNames(fullnames ...string) error {
	for _, fullname := range fullnames {
		nameSlice := strings.Split(fullname, util.SepFullFeatureName)
		if len(nameSlice) != 2 {
			return errdefs.Errorf("invalid full feature name: '%s'", fullname)
		}
	}
	return nil
}
