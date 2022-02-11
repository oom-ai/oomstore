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
	exportResult, err := store.ChannelExport(ctx, opt)
	if err != nil {
		return err
	}
	for row := range exportResult.Data {
   		if row.Error != nil {
			return err
    		}
		fmt.Println(cast.ToStringSlice([]interface{}(row.Record)))
	}
    return nil
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
	if len(features) == 0 {
		return nil, nil
	}

	if len(features) != len(opt.FeatureNames) {
		invalid := features.FindMissingFeatures(opt.FeatureNames)
		return nil, errdefs.Errorf("invalid feature names %s", invalid)
	}

	featureMap := features.GroupByGroupID()
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
			return nil, errdefs.Errorf("group %s no feature values up to %d, use a later timestamp", group.Name, opt.UnixMilli)
		} else {
			if revision.SnapshotTable == "" {
				if err = s.Snapshot(ctx, group.Name); err != nil {
					return nil, err
				}
			}

			if revision, err = s.GetRevision(ctx, revision.ID); err != nil {
				return nil, err
			}
			snapshotTables[group.ID] = revision.SnapshotTable
			if group.Category == types.CategoryStream {
				cdcTables[group.ID] = revision.CdcTable
			}
		}
	}

	return s.offline.Export(ctx, offline.ExportOpt{
		SnapshotTables: snapshotTables,
		CdcTables:      cdcTables,
		Features:       featureMap,
		UnixMilli:      opt.UnixMilli,
		EntityName:     features[0].Group.Entity.Name,
		Limit:          opt.Limit,
	})
}

// Export exports the latest feature value up to the given timestamp, it outputs
// feature values to the given file path.
func (s *OomStore) Export(ctx context.Context, opt types.ExportOpt) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

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
		if row.Error != nil {
			return row.Error
		}

		if err := w.Write(cast.ToStringSlice([]interface{}(row.Record))); err != nil {
			return err
		}
	}
	return nil
}
