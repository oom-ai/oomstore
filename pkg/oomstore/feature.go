package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Get metadata of a feature by ID.
func (s *OomStore) GetFeature(ctx context.Context, id int) (*types.Feature, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CacheGetFeature(ctx, id)
}

// Get metadata of a feature by name.
func (s *OomStore) GetFeatureByName(ctx context.Context, name string) (*types.Feature, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	return s.metadata.CacheGetFeatureByName(ctx, name)
}

// List metadata of features meeting particular criteria.
func (s *OomStore) ListFeature(ctx context.Context, opt types.ListFeatureOpt) (types.FeatureList, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	metadataOpt := metadata.ListFeatureOpt{
		FeatureNames: opt.FeatureNames,
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
	return s.metadata.CacheListFeature(ctx, metadataOpt), nil
}

// Update metadata of a feature.
func (s *OomStore) UpdateFeature(ctx context.Context, opt types.UpdateFeatureOpt) error {
	if err := s.metadata.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	feature, err := s.metadata.CacheGetFeatureByName(ctx, opt.FeatureName)
	if err != nil {
		return err
	}
	return s.metadata.UpdateFeature(ctx, metadata.UpdateFeatureOpt{
		FeatureID:      feature.ID,
		NewDescription: opt.NewDescription,
	})
}

// Create metadata of a batch feature.
func (s *OomStore) CreateBatchFeature(ctx context.Context, opt types.CreateFeatureOpt) (int, error) {
	if err := s.metadata.Refresh(); err != nil {
		return 0, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	group, err := s.metadata.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return 0, err
	}
	if group.Category != types.BatchFeatureCategory {
		return 0, fmt.Errorf("expected batch feature group, got %s feature group", group.Category)
	}

	valueType, err := s.offline.TypeTag(opt.DBValueType)
	if err != nil {
		return 0, err
	}
	return s.metadata.CreateFeature(ctx, metadata.CreateFeatureOpt{
		FeatureName: opt.FeatureName,
		GroupID:     group.ID,
		DBValueType: opt.DBValueType,
		ValueType:   valueType,
		Description: opt.Description,
	})
}
