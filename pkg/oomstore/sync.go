package oomstore

import (
	"context"
	"fmt"
	"time"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) Sync(ctx context.Context, opt types.SyncOpt) error {
	revision, err := s.GetRevision(ctx, types.GetRevisionOpt{
		GroupName:  &opt.GroupName,
		RevisionId: &opt.RevisionId,
	})
	if err != nil {
		return err
	}

	group, err := s.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return err
	}

	if group.OnlineRevisionID != nil && *group.OnlineRevisionID == revision.ID {
		return fmt.Errorf("online store already in the latest revision")
	}

	entity, err := s.GetEntity(ctx, group.EntityName)
	if err != nil {
		return err
	}

	features, err := s.ListFeature(ctx, types.ListFeatureOpt{GroupName: &opt.GroupName})
	if err != nil {
		return err
	}

	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   group.EntityName,
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

	var previousRevision *types.Revision
	if group.OnlineRevisionID != nil {
		previousRevision, err = s.metadata.GetRevision(ctx, metadata.GetRevisionOpt{
			RevisionId: group.OnlineRevisionID,
		})
		if err != nil {
			return err
		}
	}

	if _, err = s.metadata.UpdateFeatureGroup(ctx, types.UpdateFeatureGroupOpt{
		GroupName:        group.Name,
		OnlineRevisionId: &revision.ID,
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
		if _, err = s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
			RevisionID:  int64(revision.ID),
			NewRevision: &newRevision,
			NewAnchored: &newChored,
		}); err != nil {
			return err
		}
	}
	return nil
}
