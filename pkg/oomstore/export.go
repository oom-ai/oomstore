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

	var dataTable string
	if opt.GroupRevision == nil {
		revision, err := s.GetRevision(ctx, metadatav2.GetRevisionOpt{
			GroupID:  &opt.GroupID,
			Revision: opt.GroupRevision,
		})
		if err != nil {
			return nil, nil, err
		}
		dataTable = revision.DataTable
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

	entity, err := s.metadatav2.GetEntity(ctx, group.EntityID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get entity id=%d: %v", group.EntityID, err)
	}

	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    dataTable,
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
