package onestore

import (
	"context"
	"fmt"

	dbtypes "github.com/onestore-ai/onestore/internal/database/types"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) ExportFeatureValues(ctx context.Context, opt types.ExportFeatureValuesOpt) ([]string, <-chan *types.RawFeatureValueRecord, error) {
	group, err := s.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return nil, nil, err
	}

	var dataTable string
	if opt.GroupRevision == nil {
		if group.DataTable == nil {
			return nil, nil, fmt.Errorf("feature group '%s' data source not set", opt.GroupName)
		}
		dataTable = *group.DataTable
	} else {
		revision, err := s.GetRevision(ctx, opt.GroupName, *opt.GroupRevision)
		if err != nil {
			return nil, nil, err
		}
		dataTable = revision.DataTable
	}

	featureNames := opt.FeatureNames
	allFeatures, err := s.ListFeature(ctx, types.ListFeatureOpt{GroupName: &opt.GroupName})
	if err != nil {
		return nil, nil, err
	}

	allFeatureNames := make([]string, 0, len(allFeatures))
	for _, f := range allFeatures {
		allFeatureNames = append(allFeatureNames, f.Name)
	}

	if len(featureNames) == 0 {
		featureNames = allFeatureNames
	} else {
		for _, field := range featureNames {
			if !contains(allFeatureNames, field) {
				return nil, nil, fmt.Errorf("feature '%s' does not exist", field)
			}
		}
	}

	stream, err := s.offline.GetFeatureValuesStream(ctx, dbtypes.GetFeatureValuesStreamOpt{
		DataTable:    dataTable,
		EntityName:   group.EntityName,
		FeatureNames: featureNames,
		Limit:        opt.Limit,
	})

	fields := append([]string{group.EntityName}, featureNames...)
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
