package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// ListRevision lists metadata of revisions of a same group.
func (s *OomStore) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	return s.metadata.ListRevision(ctx, groupID)
}

// GetRevision gets metadata of a revision by ID.
func (s *OomStore) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	return s.metadata.GetRevision(ctx, id)
}

// GetRevisionBy gets metadata of a revision by group ID and revision.
func (s *OomStore) GetRevisionBy(ctx context.Context, groupID int, revision int64) (*types.Revision, error) {
	return s.metadata.GetRevisionBy(ctx, groupID, revision)
}

// createRevision creates a new revision without snapshot table or cdc table.
func (s *OomStore) createRevision(ctx context.Context, opt metadata.CreateRevisionOpt) (int, error) {
	if err := s.createDummyRevisionAndTables(ctx, opt.GroupID); err != nil {
		return 0, err
	}

	return s.metadata.CreateRevision(ctx, opt)
}

// createRevisionAndCdcTable creates a new revision with cdc table.
func (s *OomStore) createRevisionAndCdcTable(ctx context.Context, groupID int, revision int64) error {
	features := s.metadata.ListCachedFeature(ctx, metadata.ListCachedFeatureOpt{
		GroupID: &groupID,
	})
	entity := features[0].Entity()

	if err := s.offline.CreateTable(ctx, offline.CreateTableOpt{
		TableName:  dbutil.OfflineStreamCdcTableName(groupID, revision),
		EntityName: entity.Name,
		Features:   features,
		TableType:  types.TableStreamCdc,
	}); err != nil {
		return err
	}

	snapshotTable := ""
	cdcTable := dbutil.OfflineStreamCdcTableName(groupID, revision)
	if _, err := s.createRevision(ctx, metadata.CreateRevisionOpt{
		GroupID:       groupID,
		Revision:      revision,
		SnapshotTable: &snapshotTable,
		CdcTable:      &cdcTable,
	}); err != nil {
		return err
	}
	return nil
}

// createDummyRevisionAndTables creates dummy revision (revision = 0) with snapshot table and cdc table.
func (s *OomStore) createDummyRevisionAndTables(ctx context.Context, groupID int) error {
	_, err := s.GetRevisionBy(ctx, groupID, 0)
	if err == nil {
		return nil
	}
	if !errdefs.IsNotFound(err) {
		return err
	}

	revisionID, err := s.metadata.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:    0,
		GroupID:     groupID,
		Description: "dummy revision",
	})
	if err != nil {
		return err
	}
	return s.createSnapshotAndCdcTable(ctx, revisionID)
}

// createSnapshotAndCdcTable creates snapshot table and cdc table for a specified revision.
func (s *OomStore) createSnapshotAndCdcTable(ctx context.Context, revisionID int) error {
	revision, err := s.GetRevision(ctx, revisionID)
	if err != nil {
		return err
	}
	var snapshotTableName string
	if revision.Group.Category == types.CategoryStream {
		snapshotTableName = dbutil.OfflineStreamSnapshotTableName(revision.GroupID, revision.Revision)
	} else {
		snapshotTableName = dbutil.OfflineBatchSnapshotTableName(revision.GroupID, int64(revision.ID))
	}

	// Create snapshot table in offline store
	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &revision.GroupID,
	})
	if err != nil {
		return err
	}

	if err = s.offline.CreateTable(ctx, offline.CreateTableOpt{
		TableName:  snapshotTableName,
		EntityName: revision.Group.Entity.Name,
		Features:   features,
		TableType:  types.TableStreamSnapshot,
	}); err != nil {
		return err
	}

	var cdcTable *string
	if revision.Group.Category == types.CategoryStream {
		tableName := dbutil.OfflineStreamCdcTableName(revision.GroupID, revision.Revision)
		if err = s.offline.CreateTable(ctx, offline.CreateTableOpt{
			TableName:  tableName,
			EntityName: revision.Group.Entity.Name,
			Features:   features,
			TableType:  types.TableStreamCdc,
		}); err != nil {
			return err
		}
		cdcTable = &tableName
	}

	// Update snapshot_table in feature_group_revision table
	return s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
		RevisionID:       revision.ID,
		NewSnapshotTable: &snapshotTableName,
		NewCdcTable:      cdcTable,
	})
}
