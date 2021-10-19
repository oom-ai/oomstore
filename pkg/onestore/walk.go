package onestore

import (
	"context"
	"fmt"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) WalkFeatureValues(ctx context.Context, opt types.WalkFeatureValuesOpt) error {
	fields := opt.FeatureNames
	allFeatures, err := s.ListFeature(ctx, types.ListFeatureOpt{GroupName: &opt.FeatureGroup.Name})
	if err != nil {
		return err
	}

	allFeatureNames := make([]string, 0, len(allFeatures))
	for _, f := range allFeatures {
		allFeatureNames = append(allFeatureNames, f.Name)
	}

	if len(fields) == 0 {
		fields = allFeatureNames
	} else {
		for _, field := range fields {
			if !contains(allFeatureNames, field) {
				return fmt.Errorf("feature '%s' does not exist", field)
			}
		}
	}

	// set entity key as the first field
	fields = append([]string{opt.FeatureGroup.EntityName}, fields...)

	table := opt.FeatureGroup.DataTable
	if table == nil {
		return fmt.Errorf("feature group '%s' data source not set", opt.FeatureGroup.Name)
	}

	return s.db.WalkTable(ctx, *table, fields, opt.Limit, func(slice []interface{}) error {
		if len(slice) < 1 {
			return fmt.Errorf("empty row")
		}
		key := slice[0].(string)
		return opt.WalkFeatureValuesFunc(key, slice[1:])
	})
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
