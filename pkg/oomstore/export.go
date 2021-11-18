package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Export feature values of a particular revision.
func (s *OomStore) ExportFeatureValues(ctx context.Context, opt types.ExportFeatureValuesOpt) ([]string, <-chan *types.RawFeatureValueRecord, error) {
	revision, err := s.GetRevision(ctx, opt.RevisionID)
	if err != nil {
		return nil, nil, err
	}

	featureNames := opt.FeatureNames
	allFeatures, err := s.ListFeature(ctx, types.ListFeatureOpt{
		GroupName: &revision.Group.Name,
	})
	if err != nil {
		return nil, nil, err
	}

	allFeatureNames := allFeatures.Names()
	if len(featureNames) == 0 {
		featureNames = allFeatureNames
	} else {
		for _, field := range featureNames {
			if !contains(allFeatureNames, field) {
				return nil, nil, fmt.Errorf("feature '%s' does not exist", field)
			}
		}
	}

	entity := revision.Group.Entity
	if entity == nil {
		return nil, nil, fmt.Errorf("failed to get entity id='%d'", revision.GroupID)
	}

	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   entity.Name,
		FeatureNames: featureNames,
		Limit:        opt.Limit,
	})

	fields := append([]string{entity.Name}, featureNames...)
	return fields, stream, err
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
