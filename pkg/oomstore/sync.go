package oomstore

import (
	"context"
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

func (s *OomStore) Sync(ctx context.Context, opt types.SyncOpt) error {
	revision, err := s.GetRevision(ctx, metadatav2.GetRevisionOpt{
		GroupID:    &opt.GroupID,
		RevisionId: &opt.RevisionId,
	})
	if err != nil {
		return err
	}

	group, err := s.GetFeatureGroup(ctx, opt.GroupID)
	if err != nil {
		return err
	}

	if group.OnlineRevisionID != nil && *group.OnlineRevisionID == revision.ID {
		return fmt.Errorf("online store already in the latest revision")
	}

	entity, err := s.GetEntity(ctx, group.EntityID)
	if err != nil {
		return err
	}

	features := s.ListFeature(ctx, metadatav2.ListFeatureOpt{GroupID: &opt.GroupID})
	if err != nil {
		return err
	}

	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   entity.Name,
		FeatureNames: features.Names(),
	})
	if err != nil {
		return err
	}

	if err = s.online.Import(ctx, online.ImportOpt{
		Features: features,
		Revision: revision,
		Entity:   entity,
		Stream:   stream,
	}); err != nil {
		return err
	}

	var previousRevision *typesv2.Revision
	if group.OnlineRevisionID != nil {
		previousRevision, err = s.metadatav2.GetRevision(ctx, metadatav2.GetRevisionOpt{
			RevisionId: group.OnlineRevisionID,
		})
		if err != nil {
			return err
		}
	}

	if err = s.metadatav2.UpdateFeatureGroup(ctx, metadatav2.UpdateFeatureGroupOpt{
		GroupID:             group.ID,
		NewOnlineRevisionID: &revision.ID,
	}); err != nil {
		return err
	}

	if previousRevision != nil {
		return s.online.Purge(ctx, previousRevision)
	}

	if !revision.Anchored {
		newRevision := time.Now().Unix()
		newChored := true
		// update revision timestamp using current timestamp
		if err = s.metadatav2.UpdateRevision(ctx, metadatav2.UpdateRevisionOpt{
			RevisionID:  int32(revision.ID),
			NewRevision: &newRevision,
			NewAnchored: &newChored,
		}); err != nil {
			return err
		}
	}
	return nil
}
