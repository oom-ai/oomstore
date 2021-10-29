package oomstore

import (
	"context"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (s *OomStore) Materialize(ctx context.Context, opt types.MaterializeOpt) error {
	group, err := s.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return err
	}

	entity, err := s.GetEntity(ctx, group.EntityName)
	if err != nil {
		return err
	}

	features, err := s.ListFeature(ctx, types.ListFeatureOpt{GroupName: &opt.GroupName})
	if err != nil {
		return err
	}

	revision, err := s.getMaterializeRevision(ctx, opt)
	if err != nil {
		return err
	}
	if group.OnlineRevisionID != nil && *group.OnlineRevisionID == revision.ID {
		return fmt.Errorf("online store already in the latest revision")
	}

	stream, err := s.offline.Export(ctx, offline.ExportOpt{
		DataTable:    revision.DataTable,
		EntityName:   group.EntityName,
		FeatureNames: features.Names(),
	})
	if err != nil {
		return err
	}

	err = s.online.Import(ctx, online.ImportOpt{
		Features: features,
		Revision: revision,
		Entity:   entity,
		Stream:   stream,
	})
	if err != nil {
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
		return s.online.Purge(ctx, revision)
	}
	return nil
}

func (s *OomStore) getMaterializeRevision(ctx context.Context, opt types.MaterializeOpt) (*types.Revision, error) {
	if opt.GroupRevision != nil {
		return s.GetRevision(ctx, opt.GroupName, *opt.GroupRevision)
	}
	return s.metadata.GetLatestRevision(ctx, opt.GroupName)
}
