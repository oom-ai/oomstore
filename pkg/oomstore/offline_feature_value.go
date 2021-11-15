package oomstore

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/oom-ai/oomstore/internal/database/metadatav2"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/typesv2"
)

// GetHistoricalFeatureValues gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OomStore) GetHistoricalFeatureValues(ctx context.Context, opt types.GetHistoricalFeatureValuesOpt) (*types.JoinResult, error) {
	features := s.metadatav2.ListFeature(ctx, metadatav2.ListFeatureOpt{FeatureIDs: &opt.FeatureIDs})

	features = features.Filter(func(f *typesv2.Feature) bool {
		return f.Group.Category == types.BatchFeatureCategory
	})
	if len(features) == 0 {
		return nil, nil
	}

	entity, err := s.getSharedEntity(features)
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, fmt.Errorf("failed to get shared entity")
	}

	featureMap := buildGroupToFeaturesMap(features)
	revisionRangeMap := make(map[string][]*metadatav2.RevisionRange)
	for groupName := range featureMap {
		// TODO: This is slow but I haven't figured out a better way
		group, err := s.metadatav2.GetFeatureGroupByName(ctx, groupName)
		if err != nil {
			return nil, err
		}
		revisionRanges, err := s.buildRevisionRanges(ctx, group.ID)
		if err != nil {
			return nil, err
		}
		revisionRangeMap[groupName] = revisionRanges
	}

	return s.offline.Join(ctx, offline.JoinOpt{
		Entity:           *entity,
		EntityRows:       opt.EntityRows,
		FeatureMap:       featureMap,
		RevisionRangeMap: revisionRangeMap,
	})
}

// key: group_name, value: slice of features
func buildGroupToFeaturesMap(features typesv2.FeatureList) map[string]typesv2.FeatureList {
	groups := make(map[string]typesv2.FeatureList)
	for _, f := range features {
		if _, ok := groups[f.Group.Name]; !ok {
			groups[f.Group.Name] = typesv2.FeatureList{}
		}
		groups[f.Group.Name] = append(groups[f.Group.Name], f)
	}
	return groups
}

func (s *OomStore) buildRevisionRanges(ctx context.Context, groupID int16) ([]*metadatav2.RevisionRange, error) {
	revisions := s.metadatav2.ListRevision(ctx, metadatav2.ListRevisionOpt{GroupID: &groupID})
	if len(revisions) == 0 {
		return nil, nil
	}

	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})

	var ranges []*metadatav2.RevisionRange
	for i := 1; i < len(revisions); i++ {
		ranges = append(ranges, &metadatav2.RevisionRange{
			MinRevision: revisions[i-1].Revision,
			MaxRevision: revisions[i].Revision,
			DataTable:   revisions[i-1].DataTable,
		})
	}

	return append(ranges, &metadatav2.RevisionRange{
		MinRevision: revisions[len(revisions)-1].Revision,
		MaxRevision: math.MaxInt64,
		DataTable:   revisions[len(revisions)-1].DataTable,
	}), nil
}
