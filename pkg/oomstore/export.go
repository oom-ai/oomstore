package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

/*
Export feature values of a particular revision.
Usage Example:
	exportResult, err := store.Export(ctx, opt)
	if err != nil {
		return err
	}
	for row := range exportResult.Data {
		fmt.Println(cast.ToStringSlice([]interface{}(row)))
	}
	// Attention: call CheckStreamError after consuming exportResult.Data channel
	return exportResult.CheckStreamError()
*/
func (s *OomStore) Export(ctx context.Context, opt types.ExportOpt) (*types.ExportResult, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	revision, err := s.GetRevision(ctx, opt.RevisionID)
	if err != nil {
		return nil, err
	}

	featureNames := opt.FeatureNames
	allFeatures, err := s.ListFeature(ctx, types.ListFeatureOpt{
		GroupName: &revision.Group.Name,
	})
	if err != nil {
		return nil, err
	}

	allFeatureNames := allFeatures.Names()
	if len(featureNames) == 0 {
		featureNames = allFeatureNames
	} else {
		for _, field := range featureNames {
			if !contains(allFeatureNames, field) {
				return nil, fmt.Errorf("feature '%s' does not exist", field)
			}
		}
	}

	entity := revision.Group.Entity
	if entity == nil {
		return nil, fmt.Errorf("failed to get entity id=%d", revision.GroupID)
	}

	stream, exportErr := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   entity.Name,
		FeatureNames: featureNames,
		Limit:        opt.Limit,
	})
	header := append([]string{entity.Name}, featureNames...)
	return types.NewExportResult(header, stream, exportErr), nil
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
