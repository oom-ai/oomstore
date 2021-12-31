package oomstore

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/offline"

	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"

	"github.com/oom-ai/oomstore/internal/database/metadata"

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

func (s *OomStore) createFirstSnapshotTable(ctx context.Context, revision *types.Revision) error {
	snapshotTable := sqlutil.OfflineStreamSnapshotTableName(revision.GroupID, revision.Revision)

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
		TableName: snapshotTable,
		Entity:    revision.Group.Entity,
		Features:  features,
	}); err != nil {
		return err
	}

	return nil
}
