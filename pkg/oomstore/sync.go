package oomstore

import (
	"context"
	"fmt"
	"time"

	"github.com/pkg/errors"

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
	if err := s.metadata.Refresh(); err != nil {
		return fmt.Errorf("failed to refresh informer, err=%v", err)
	}
	revision, err := s.GetRevision(ctx, opt.RevisionID)
	if err != nil {
		return err
	}

	group := revision.Group
	prevOnlineRevisionID := group.OnlineRevisionID
	if prevOnlineRevisionID != nil && *prevOnlineRevisionID == opt.RevisionID {
		return errors.Errorf("the specific revision was synced to the online store, won't do it again this time")
	}

	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		GroupName: &group.Name,
	})
	if err != nil {
		return err
	}

	// Move data from offline to online store
	exportStream, exportError := s.offline.Export(ctx, offline.ExportOpt{
		SnapshotTable: revision.SnapshotTable,
		EntityName:    group.Entity.Name,
		Features:      features,
	})

	if err = s.online.Import(ctx, online.ImportOpt{
		Features:     features,
		Revision:     revision,
		Entity:       group.Entity,
		ExportStream: exportStream,
		ExportError:  exportError,
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
			newRevision := time.Now().Unix()
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
