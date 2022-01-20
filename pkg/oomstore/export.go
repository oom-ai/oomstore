package oomstore

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/util"
)

/*
ChannelExport exports the latest feature value up to the given timestamp.
Usage Example:
	exportResult, err := store.ExportBatch(ctx, opt)
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
	if err := util.ValidateFullFeatureNames(opt.FeatureNames...); err != nil {
		return nil, err
	}
	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		FeatureNames: &opt.FeatureNames,
	})
	if err != nil {
		return nil, err
	}
	if len(features) != len(opt.FeatureNames) {
		invalid := make([]string, 0)
		for _, name := range opt.FeatureNames {
			f := features.Find(func(feature *types.Feature) bool {
				return feature.FullName() == name
			})
			if f == nil {
				invalid = append(invalid, name)
			}
		}
		return nil, errdefs.Errorf("invalid feature names %s", invalid)
	}

	featureMap := buildGroupIDToFeaturesMap(features)
	snapshotTables := make(map[int]string)
	cdcTables := make(map[int]string)
	for _, featureList := range featureMap {
		group := featureList[0].Group
		revisions, err := s.ListRevision(ctx, &group.ID)
		if err != nil {
			return nil, err
		}
		revision := revisions.Before(opt.UnixMilli)
		if revision == nil {
			empty := make(chan types.ExportRecord)
			close(empty)
			return &types.ExportResult{
				Header: append([]string{features[0].Entity().Name}, features.FullNames()...),
				Data:   empty,
			}, nil
		}
		if revision.SnapshotTable == "" {
			if err = s.Snapshot(ctx, group.Name); err != nil {
				return nil, err
			}
		}
		revision, err = s.GetRevision(ctx, revision.ID)
		if err != nil {
			return nil, err
		}
		snapshotTables[group.ID] = revision.SnapshotTable
		if group.Category == types.CategoryStream {
			cdcTables[group.ID] = revision.CdcTable
		}
	}

	result, err := s.offline.Export(ctx, offline.ExportOpt{
		SnapshotTables: snapshotTables,
		CdcTables:      cdcTables,
		Features:       featureMap,
		UnixMilli:      opt.UnixMilli,
		EntityName:     features[0].Group.Entity.Name,
		Limit:          opt.Limit,
	})
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Export exports the latest feature value up to the given timestamp, it outputs
// feature values to the given file path.
func (s *OomStore) Export(ctx context.Context, opt types.ExportOpt) error {
	exportResult, err := s.ChannelExport(ctx, types.ChannelExportOpt{
		FeatureNames: opt.FeatureNames,
		UnixMilli:    opt.UnixMilli,
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

// key: group_Id, value: slice of features
func buildGroupIDToFeaturesMap(features types.FeatureList) map[int]types.FeatureList {
	groups := make(map[int]types.FeatureList)
	for _, f := range features {
		if _, ok := groups[f.Group.ID]; !ok {
			groups[f.Group.ID] = types.FeatureList{}
		}
		groups[f.Group.ID] = append(groups[f.Group.ID], f)
	}
	return groups
}
