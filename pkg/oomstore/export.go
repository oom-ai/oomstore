package oomstore

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

/*
ChannelExportBatch exports batch feature values of a particular revision.
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
func (s *OomStore) ChannelExportBatch(ctx context.Context, opt types.ChannelExportBatchOpt) (*types.ExportResult, error) {
	revision, err := s.GetRevision(ctx, opt.RevisionID)
	if err != nil {
		return nil, err
	}

	var features types.FeatureList
	if len(opt.FeatureNames) == 0 {
		features, err = s.ListFeature(ctx, types.ListFeatureOpt{
			GroupName: &revision.Group.Name,
		})
	} else {
		fullNames := make([]string, 0, len(opt.FeatureNames))
		for _, name := range opt.FeatureNames {
			fullNames = append(fullNames, fmt.Sprintf("%s.%s", revision.Group.Name, name))
		}
		features, err = s.ListFeature(ctx, types.ListFeatureOpt{
			FeatureFullNames: &fullNames,
		})
	}
	if err != nil {
		return nil, err
	}

	entity := revision.Group.Entity
	if entity == nil {
		return nil, errdefs.Errorf("failed to get entity id=%d", revision.GroupID)
	}

	stream, exportErr := s.offline.ExportOneGroup(ctx, offline.ExportOneGroupOpt{
		SnapshotTable: revision.SnapshotTable,
		EntityName:    entity.Name,
		Features:      features,
		Limit:         opt.Limit,
	})
	header := append([]string{entity.Name}, features.Names()...)
	return types.NewExportResult(header, stream, exportErr), nil
}

func (s *OomStore) ExportBatch(ctx context.Context, opt types.ExportBatchOpt) error {
	exportResult, err := s.ChannelExportBatch(ctx, types.ChannelExportBatchOpt{
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

// ChannelExportStream exports the latest streaming feature values up to the given timestamp.
// Currently, this API can only export features in one feature group.
func (s *OomStore) ChannelExportStream(ctx context.Context, opt types.ChannelExportStreamOpt) (*types.ExportResult, error) {
	if err := validateFeatureFullNames(opt.FeatureFullNames); err != nil {
		return nil, errdefs.WithStack(err)
	}
	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		FeatureFullNames: &opt.FeatureFullNames,
	})
	if err != nil {
		return nil, err
	}
	if len(features.GroupIDs()) != 1 {
		return nil, errdefs.Errorf("expected 1 group, got %d groups", len(features.GroupIDs()))
	}
	group := features[0].Group
	revisions, err := s.ListRevision(ctx, &group.ID)
	if err != nil {
		return nil, err
	}
	revision := revisions.Before(opt.UnixMilli)
	if revision == nil {
		return nil, errdefs.Errorf("no feature values up to %d, use a later timestamp", opt.UnixMilli)
	}
	if revision.SnapshotTable == "" {
		if err = s.Snapshot(ctx, group.Name); err != nil {
			return nil, errdefs.WithStack(err)
		}
	}

	snapshotTable := dbutil.OfflineStreamSnapshotTableName(group.ID, revision.Revision)
	cdcTable := dbutil.OfflineStreamCdcTableName(group.ID, revision.Revision)
	stream, exportErr := s.offline.ExportOneGroup(ctx, offline.ExportOneGroupOpt{
		SnapshotTable: snapshotTable,
		CdcTable:      &cdcTable,
		UnixMilli:     &opt.UnixMilli,
		EntityName:    group.Entity.Name,
		Features:      features,
		Limit:         opt.Limit,
	})

	header := append([]string{group.Entity.Name}, features.Names()...)
	return types.NewExportResult(header, stream, exportErr), nil
}

// ChannelExport exports the latest streaming feature values up to the given timestamp.
// Currently, this API can only export features in one feature group.
func (s *OomStore) ChannelExport(ctx context.Context, opt types.ChannelExportOpt) (*types.ExportResult, error) {
	if err := validateFeatureFullNames(opt.FeatureFullNames); err != nil {
		return nil, err
	}
	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		FeatureFullNames: &opt.FeatureFullNames,
	})
	if err != nil {
		return nil, err
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
			return nil, errdefs.Errorf("no feature values up to %d, use a later timestamp", opt.UnixMilli)
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

	stream, exportErr := s.offline.Export(ctx, offline.ExportOpt{
		SnapshotTables: snapshotTables,
		CdcTables:      cdcTables,
		Features:       featureMap,
		UnixMilli:      opt.UnixMilli,
		EntityName:     features[0].Group.Entity.Name,
		Limit:          opt.Limit,
	})

	header := append([]string{features[0].Group.Entity.Name}, features.Names()...)
	return types.NewExportResult(header, stream, exportErr), nil
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
