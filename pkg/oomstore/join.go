package oomstore

import (
	"context"
	"fmt"
	"math"
	"sort"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Join gets point-in-time feature values for each entity row;
// currently, this API only supports batch features.
func (s *OomStore) Join(ctx context.Context, opt types.JoinOpt) (*types.JoinResult, error) {
	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{FeatureIDs: &opt.FeatureIDs})

	features = features.Filter(func(f *types.Feature) bool {
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
	revisionRangeMap := make(map[string][]*metadata.RevisionRange)
	for groupName, featureList := range featureMap {
		if len(featureList) == 0 {
			continue
		}
		revisionRanges, err := s.buildRevisionRanges(ctx, featureList[0].GroupID)
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
func buildGroupToFeaturesMap(features types.FeatureList) map[string]types.FeatureList {
	groups := make(map[string]types.FeatureList)
	for _, f := range features {
		if _, ok := groups[f.Group.Name]; !ok {
			groups[f.Group.Name] = types.FeatureList{}
		}
		groups[f.Group.Name] = append(groups[f.Group.Name], f)
	}
	return groups
}

func (s *OomStore) buildRevisionRanges(ctx context.Context, groupID int) ([]*metadata.RevisionRange, error) {
	revisions := s.metadata.ListRevision(ctx, metadata.ListRevisionOpt{GroupID: &groupID})
	if len(revisions) == 0 {
		return nil, nil
	}

	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})

	var ranges []*metadata.RevisionRange
	for i := 1; i < len(revisions); i++ {
		ranges = append(ranges, &metadata.RevisionRange{
			MinRevision: revisions[i-1].Revision,
			MaxRevision: revisions[i].Revision,
			DataTable:   revisions[i-1].DataTable,
		})
	}

	return append(ranges, &metadata.RevisionRange{
		MinRevision: revisions[len(revisions)-1].Revision,
		MaxRevision: math.MaxInt64,
		DataTable:   revisions[len(revisions)-1].DataTable,
	}), nil
}
