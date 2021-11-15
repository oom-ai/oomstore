package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) ExportFeatureValues(ctx context.Context, opt types.ExportFeatureValuesOpt) ([]string, <-chan *types.RawFeatureValueRecord, error) {
	group, err := s.metadatav2.GetFeatureGroup(ctx, opt.GroupID)
	if err != nil {
		return nil, nil, err
	}

	revision, err := s.GetRevision(ctx, opt.RevisionID)
	if err != nil {
		return nil, nil, err
	}

	featureNames := opt.FeatureNames
	allFeatures := s.ListFeature(ctx, metadatav2.ListFeatureOpt{GroupID: &opt.GroupID})
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

	if group.Entity == nil {
		return nil, nil, fmt.Errorf("failed to get entity id='%d'", group.ID)
	}

	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   group.Entity.Name,
		FeatureNames: featureNames,
		Limit:        opt.Limit,
	})

	fields := append([]string{group.Entity.Name}, featureNames...)
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
