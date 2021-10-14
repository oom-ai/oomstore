package onestore

import (
	"context"
	"fmt"

	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) WalkFeatureValues(ctx context.Context, opt types.WalkFeatureValuesOpt) error {
	fields := opt.FeatureNames
	// TODO: validate field name
	if len(fields) == 0 {
		fields = []string{"*"}
	} else {
		entity := opt.FeatureGroup.EntityName
		if contains(fields, entity) {
			return fmt.Errorf("%s is not a feature", entity)
		}

		fields = append([]string{entity}, fields...)
	}

	table := opt.FeatureGroup.DataTable
	return s.db.WalkTable(ctx, table, fields, opt.Limit, opt.WalkFeatureValuesFunc)
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
