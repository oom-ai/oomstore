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

	stream, exportErr := s.offline.Export(ctx, offline.ExportOpt{
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
		return nil, errdefs.WithStack(err)
	}
	if len(features.GroupIDs()) != 1 {
		return nil, fmt.Errorf("expected 1 group, got %d groups", len(features.GroupIDs()))
	}
	group := features[0].Group
	revisions, err := s.ListRevision(ctx, &group.ID)
	if err != nil {
		return nil, errdefs.WithStack(err)
	}
	revision := revisions.Before(opt.UnixMilli)
	if revision == nil {
		return nil, fmt.Errorf("no feature values up to %d, use a later timestamp", opt.UnixMilli)
	}
	if revision.SnapshotTable == "" {
		if err = s.Snapshot(ctx, group.Name); err != nil {
			return nil, errdefs.WithStack(err)
		}
	}

	snapshotTable := dbutil.OfflineStreamSnapshotTableName(group.ID, revision.Revision)
	cdcTable := dbutil.OfflineStreamCdcTableName(group.ID, revision.Revision)
	stream, exportErr := s.offline.Export(ctx, offline.ExportOpt{
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
