package oomstore

import (
	"context"
	"fmt"
	"strings"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	FeatureFullNameSeparator = "."
)

// Get metadata of a feature by ID.
func (s *OomStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	return s.metadata.GetFeature(ctx, id)
}

// Get metadata of a feature by full name.
func (s *OomStore) GetFeatureByName(ctx context.Context, fullName string) (*types.Feature, error) {
	return s.metadata.GetFeatureByName(ctx, fullName)
}

// List metadata of features meeting particular criteria.
func (s *OomStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	metadataOpt := metadata.ListFeatureOpt{
		FeatureFullNames: opt.FeatureFullNames,
	}
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
	return s.metadata.ListFeature(ctx, metadataOpt)
}

// Update metadata of a feature.
func (s *OomStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	feature, err := s.metadata.GetFeatureByName(ctx, opt.FeatureFullName)
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
		FullName:    fmt.Sprintf("%s.%s", group.Name, opt.FeatureName),
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

func validateFeatureFullNames(names []string) error {
	for _, name := range names {
		nameSlice := strings.Split(name, FeatureFullNameSeparator)
		if len(nameSlice) != 2 {
			return fmt.Errorf("invalid feature full name %s", name)
		}
	}
	return nil
}
