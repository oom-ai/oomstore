package oomstore

import (
	"context"
	"sort"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
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
	if len(revisions) == 0 {
		return nil
	}
	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})
	if revisions[0].SnapshotTable == "" {
		if err = s.createFirstSnapshotTable(ctx, revisions[0]); err != nil {
			return err
		}
	}
	for i, revision := range revisions {
		if revision.SnapshotTable != "" {
			continue
		}
		tableName := dbutil.OfflineStreamSnapshotTableName(group.ID, revision.Revision)
		if err = s.offline.Snapshot(ctx, offline.SnapshotOpt{
			Group:        group,
			Revision:     revisions[i].Revision,
			PrevRevision: revisions[i-1].Revision,
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
