package onestore

import (
	"context"

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

	revision, err := s.GetRevision(ctx, opt.GroupName, opt.GroupRevision)
	if err != nil {
		return err
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

	if err = s.metadata.UpdateFeatureGroupRevision(ctx, revision.Revision, revision.DataTable, revision.GroupName); err != nil {
		return err
	}

	return s.online.DeprecateFeatureValues(ctx, revision.GetOnlineBatchTableName())
}
