package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// List metadata of revisions of a same group.
func (s *OomStore) ListRevision(ctx context.Context, groupID *int) (types.RevisionList, error) {
	return s.metadata.ListRevision(ctx, groupID)
}

// Get metadata of a revision by ID.
func (s *OomStore) GetRevision(ctx context.Context, id int) (*types.Revision, error) {
	return s.metadata.GetRevision(ctx, id)
}

// Get metadata of a revision by group ID and revision.
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

func (s *OomStore) createFirstSnapshotTable(ctx context.Context, revision *types.Revision) error {
	snapshotTable := dbutil.OfflineStreamSnapshotTableName(revision.GroupID, revision.Revision)

	// Update snapshot_table in feature_group_revision table
	err := s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
		RevisionID:       revision.ID,
		NewSnapshotTable: &snapshotTable,
	})
	if err != nil {
		return err
	}

	// Create snapshot table in offline store
	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &revision.GroupID,
	})
	if err != nil {
		return err
	}
	if err = s.offline.CreateTable(ctx, offline.CreateTableOpt{
		TableName:  snapshotTable,
		EntityName: revision.Group.Entity.Name,
		Features:   features,
		TableType:  types.TableStreamSnapshot,
	}); err != nil {
		return err
	}

	return nil
}
