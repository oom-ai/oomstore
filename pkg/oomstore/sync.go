package oomstore

import (
	"context"
	"sort"
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Sync a particular revision of a feature group from offline to online store.
// It is a streaming process - it writes to online store while reading from offline store.
// This helps get rid of unwanted out-of-memory errors,
// where size of the particular revision outgrows memory limit of your machine.
func (s *OomStore) Sync(ctx context.Context, opt types.SyncOpt) error {
	group, revision, err := s.validateSyncOpt(ctx, opt)
	if err != nil {
		return err
	}
	if group.Category == types.CategoryStream {
		return s.syncStream(ctx, opt)
	}
	return s.syncBatch(ctx, opt, group, revision)
}

// syncBatch syncs batch feature group from offline store to online store.
func (s *OomStore) syncBatch(ctx context.Context, opt types.SyncOpt, group *types.Group, revision *types.Revision) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	prevOnlineRevisionID := group.OnlineRevisionID
	if prevOnlineRevisionID != nil && *prevOnlineRevisionID == revision.ID {
		return errdefs.Errorf("the specific revision was synced to the online store, won't do it again this time")
	}

	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		GroupNames: &[]string{group.Name},
	})
	if err != nil {
		return err
	}

	// Move data from offline to online store
	exportResult, err := s.offline.Export(ctx, offline.ExportOpt{
		SnapshotTables: map[int]string{group.ID: revision.SnapshotTable},
		Features:       map[int]types.FeatureList{group.ID: features},
		EntityName:     group.Entity.Name,
	})
	if err != nil {
		return err
	}

	if err = s.online.Import(ctx, online.ImportOpt{
		Group:        *group,
		Features:     features,
		ExportStream: exportResult.Data,
		RevisionID:   &revision.ID,
	}); err != nil {
		return err
	}

	if err = s.metadata.WithTransaction(ctx, func(c context.Context, tx metadata.DBStore) error {
		// Update the online revision id of the feature group upon sync success
		if err := tx.UpdateGroup(c, metadata.UpdateGroupOpt{
			GroupID:             group.ID,
			NewOnlineRevisionID: &revision.ID,
		}); err != nil {
			return err
		}
		if !revision.Anchored {
			newRevision := time.Now().UnixMilli()
			newChored := true
			// Update revision timestamp using current timestamp
			if err = tx.UpdateRevision(c, metadata.UpdateRevisionOpt{
				RevisionID:  revision.ID,
				NewRevision: &newRevision,
				NewAnchored: &newChored,
			}); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return err
	}

	// Now we can delete the online data corresponding to the previous revision
	if prevOnlineRevisionID != nil {
		if opt.PurgeDelay > 0 {
			time.Sleep(time.Duration(opt.PurgeDelay) * time.Second)
		}
		return s.online.Purge(ctx, *prevOnlineRevisionID)
	}
	return nil
}

// syncStream syncs stream feature group from offline store to online store.
func (s *OomStore) syncStream(ctx context.Context, opt types.SyncOpt) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		GroupNames: &[]string{opt.GroupName},
	})
	if err != nil {
		return err
	}
	group := features[0].Group
	exportResult, err := s.ChannelExport(ctx, types.ChannelExportOpt{
		FeatureNames: features.FullNames(),
		UnixMilli:    time.Now().UnixMilli(),
	})
	if err != nil {
		return err
	}

	return s.online.Import(ctx, online.ImportOpt{
		Group:        *group,
		Features:     features,
		ExportStream: exportResult.Data,
	})
}

func (s *OomStore) validateSyncOpt(ctx context.Context, opt types.SyncOpt) (*types.Group, *types.Revision, error) {
	group, err := s.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return nil, nil, err
	}
	if group.Category == types.CategoryStream && opt.RevisionID != nil {
		return nil, nil, errdefs.Errorf("streaming feature group only sync the latest values, cannot designate revisionID")
	}
	var revision *types.Revision
	if opt.RevisionID != nil {
		r, err := s.GetRevision(ctx, *opt.RevisionID)
		if err != nil {
			return nil, nil, err
		}
		revision = r
	} else {
		revisions, err := s.ListRevision(ctx, &group.ID)
		if err != nil {
			return nil, nil, err
		}
		if len(revisions) == 0 {
			return nil, nil, errdefs.Errorf("group %s doesn't have any revision", opt.GroupName)
		}
		sort.Slice(revisions, func(i, j int) bool {
			return revisions[i].Revision > revisions[j].Revision
		})
		revision = revisions[0]
	}
	if group.ID != revision.GroupID {
		return nil, nil, errdefs.Errorf("revisionID %d does not belong to group %s", *opt.RevisionID, opt.GroupName)
	}
	return group, revision, nil
}
