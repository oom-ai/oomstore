package oomstore

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
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
func (s *OomStore) ChannelExport(ctx context.Context, opt types.ChannelExportOpt) (*types.ExportResult, error) {
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

func (s *OomStore) Export(ctx context.Context, opt types.ExportOpt) error {
	exportResult, err := s.ChannelExport(ctx, types.ChannelExportOpt{
		RevisionID:   opt.RevisionID,
		FeatureNames: opt.FeatureNames,
		Limit:        opt.Limit,
	})
	if err != nil {
		return err
	}

	file, err := os.Create(opt.OutputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	if err := w.Write(exportResult.Header); err != nil {
		return err
	}
	for row := range exportResult.Data {
		if err := w.Write(cast.ToStringSlice([]interface{}(row))); err != nil {
			return err
		}
	}
	return exportResult.CheckStreamError()
}
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
