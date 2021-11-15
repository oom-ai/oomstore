package oomstore

import (
	"context"
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Move a certain feature group revision data from offline to online store
func (s *OomStore) Sync(ctx context.Context, opt types.SyncOpt) error {
	revision, err := s.GetRevision(ctx, opt.RevisionId)
	if err != nil {
		return err
	}

	group := revision.Group
	prevOnlineRevisionID := group.OnlineRevisionID
	if prevOnlineRevisionID != nil && *prevOnlineRevisionID == opt.RevisionId {
		return fmt.Errorf("the specific revision was synced to the online store, won't do it again this time")
	}

	features := s.ListFeature(ctx, metadatav2.ListFeatureOpt{GroupID: &group.ID})
	if err != nil {
		return err
	}

	// Move data from offline to online store
	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   group.Entity.Name,
		FeatureNames: features.Names(),
	})
	if err != nil {
		return err
	}

	if err = s.online.Import(ctx, online.ImportOpt{
		FeatureList: features,
		Revision:    revision,
		Entity:      group.Entity,
		Stream:      stream,
	}); err != nil {
		return err
	}

	// Update the online revision id of the feature group upon sync success
	if err = s.metadatav2.UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
		GroupID:             group.ID,
		NewOnlineRevisionID: &revision.ID,
	}); err != nil {
		return err
	}

	// Now we can delete the online data corresponding to the previous revision
	if prevOnlineRevisionID != nil {
		return s.online.Purge(ctx, *prevOnlineRevisionID)
	}

	if !revision.Anchored {
		newRevision := time.Now().Unix()
		newChored := true
		// update revision timestamp using current timestamp
		if err = s.metadatav2.UpdateRevision(ctx, metadatav2.UpdateRevisionOpt{
			RevisionID:  revision.ID,
			NewRevision: &newRevision,
			NewAnchored: &newChored,
		}); err != nil {
			return err
		}
	}
	return nil
}
