package onestore

import (
	"context"
	"fmt"

	dbtypes "github.com/onestore-ai/onestore/internal/database/types"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (s *OneStore) Materialize(ctx context.Context, opt types.MaterializeOpt) error {
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

	featureNames := []string{}
	for _, f := range features {
		featureNames = append(featureNames, f.Name)
	}

	revision, err := s.getMaterializeRevision(ctx, opt)
	if err != nil {
		return err
	}
	if group.Revision != nil && *group.Revision == revision.Revision {
		return fmt.Errorf("online store already in the latest revision")
	}

	stream, err := s.offline.GetFeatureValuesStream(ctx, dbtypes.GetFeatureValuesStreamOpt{
		DataTable:    revision.DataTable,
		EntityName:   group.EntityName,
		FeatureNames: featureNames,
	})
	if err != nil {
		return err
	}

	if err = s.online.SinkFeatureValuesStream(ctx, stream, features, revision, entity); err != nil {
		return err
	}

	var previousRevision *types.Revision
	if group.Revision != nil {
		previousRevision, err = s.GetRevision(ctx, group.Name, *group.Revision)
		if err != nil {
			return err
		}
	}

	if err = s.metadata.UpdateFeatureGroupRevision(ctx, revision.Revision, revision.DataTable, revision.GroupName); err != nil {
		return err
	}

	if previousRevision != nil {
		return s.online.PurgeRevision(ctx, revision)
	}
	return nil
}

func (s *OneStore) getMaterializeRevision(ctx context.Context, opt types.MaterializeOpt) (*types.Revision, error) {
	if opt.GroupRevision != nil {
		return s.GetRevision(ctx, opt.GroupName, *opt.GroupRevision)
	}
	return s.metadata.GetLatestRevision(ctx, opt.GroupName)
}
