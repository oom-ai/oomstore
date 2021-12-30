package oomstore

import (
	"context"
	"sort"

	"github.com/oom-ai/oomstore/internal/database/metadata/sqlutil"

	"github.com/oom-ai/oomstore/internal/database/metadata"

	"github.com/oom-ai/oomstore/internal/database/offline"
)

func (s *OomStore) Snapshot(ctx context.Context, groupName string) error {
	group, err := s.metadata.GetGroupByName(ctx, groupName)
	if err != nil {
		return err
	}
	revisions, err := s.metadata.ListRevision(ctx, &group.ID)
	if err != nil {
		return err
	}
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})
	for i, revision := range revisions {
		if revision.SnapshotTable != "" {
			continue
		}
		tableName := sqlutil.OfflineStreamSnapshotTableName(group.ID, revision.ID)
		if err = s.offline.Snapshot(ctx, offline.SnapshotOpt{
			Group:          group,
			RevisionID:     revisions[i].ID,
			PrevRevisionID: revisions[i-1].ID,
		}); err != nil {
			return err
		}
		if err = s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
			RevisionID:       revision.ID,
			NewSnapshotTable: &tableName,
		}); err != nil {
			return err
		}
	}
	return nil
}
