package oomstore

import (
	"context"
	"encoding/csv"
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/oom-ai/oomstore/pkg/oomstore/util"
)

// ChannelJoin gets point-in-time correct feature values for each entity row.
func (s *OomStore) ChannelJoin(ctx context.Context, opt types.ChannelJoinOpt) (*types.JoinResult, error) {
	if err := util.ValidateFullFeatureNames(opt.JoinFeatureNames...); err != nil {
		return nil, err
	}

	features, err := s.ListFeature(ctx, types.ListFeatureOpt{
		FeatureNames: &opt.JoinFeatureNames,
	})
	if err != nil {
		return nil, err
	}
	if len(features) == 0 {
		data := make(chan types.JoinRecord)
		close(data)

		return &types.JoinResult{Data: data}, nil
	}

	entity, err := features.GetSharedEntity()
	if err != nil {
		return nil, err
	}
	if entity == nil {
		return nil, errdefs.Errorf("failed to get shared entity")
	}

	groupNames := features.GroupNames()
	featureMap := features.GroupByGroupName()
	revisionRangeMap := make(map[string][]*offline.RevisionRange)
	for groupName, featureList := range featureMap {
		if len(featureList) == 0 {
			continue
		}
		revisionRanges, err := s.buildRevisionRanges(ctx, featureList[0].Group)
		if err != nil {
			return nil, err
		}
		if len(revisionRanges) == 0 {
			return nil, errdefs.Errorf("group %s no feature values", groupName)
		}
		revisionRangeMap[groupName] = revisionRanges
	}

	return s.offline.Join(ctx, offline.JoinOpt{
		EntityName:       entity.Name,
		EntityRows:       opt.EntityRows,
		GroupNames:       groupNames,
		FeatureMap:       featureMap,
		RevisionRangeMap: revisionRangeMap,
		ValueNames:       opt.ExistedFeatureNames,
	})
}

// Join gets point-in-time correct feature values for each entity row.
// The method is similar to ChannelJoin, except a input files on disk.
// Input File should contain header, the first two columns of Input File should be
// entity_key, unix_milli, then followed by other real-time feature values.
func (s *OomStore) Join(ctx context.Context, opt types.JoinOpt) (*types.JoinResult, error) {
	if err := util.ValidateFullFeatureNames(opt.FeatureNames...); err != nil {
		return nil, err
	}

	entityRows, header, err := getEntityRowsFromInputFile(ctx, opt.InputFilePath)
	if err != nil {
		return nil, err
	}

	return s.ChannelJoin(ctx, types.ChannelJoinOpt{
		JoinFeatureNames:    opt.FeatureNames,
		EntityRows:          entityRows,
		ExistedFeatureNames: header[2:],
	})
}

func (s *OomStore) buildRevisionRanges(ctx context.Context, group *types.Group) ([]*offline.RevisionRange, error) {
	revisions, err := s.metadata.ListRevision(ctx, &group.ID)
	if err != nil {
		return nil, err
	}
	if len(revisions) == 0 {
		return nil, nil
	}

	sort.Slice(revisions, func(i, j int) bool {
		return revisions[i].Revision < revisions[j].Revision
	})
	for _, revision := range revisions {
		if revision.SnapshotTable == "" {
			if err = s.Snapshot(ctx, group.Name); err != nil {
				return nil, err
			}
		}
	}
	revisions, err = s.metadata.ListRevision(ctx, &group.ID)
	if err != nil {
		return nil, err
	}

	var ranges []*offline.RevisionRange
	for i := 1; i < len(revisions); i++ {
		ranges = append(ranges, &offline.RevisionRange{
			MinRevision:   revisions[i-1].Revision,
			MaxRevision:   revisions[i].Revision,
			SnapshotTable: revisions[i-1].SnapshotTable,
			CdcTable:      revisions[i-1].CdcTable,
		})
	}
	ranges = append(ranges, &offline.RevisionRange{
		MinRevision:   revisions[len(revisions)-1].Revision,
		MaxRevision:   math.MaxInt64,
		SnapshotTable: revisions[len(revisions)-1].SnapshotTable,
		CdcTable:      revisions[len(revisions)-1].CdcTable,
	})
	return ranges, nil
}

func getEntityRowsFromInputFile(ctx context.Context, inputFilePath string) (<-chan types.EntityRow, []string, error) {
	input, err := os.Open(inputFilePath)
	if err != nil {
		return nil, nil, errdefs.WithStack(err)
	}
	reader := csv.NewReader(input)
	header, err := reader.Read()
	if err != nil {
		return nil, nil, errdefs.WithStack(err)
	}

	entityRows := make(chan types.EntityRow)
	go func() {
		defer close(entityRows)
		defer input.Close()

		for i := 1; ; i++ {
			line, err := reader.Read()
			if err != nil {
				if err == io.EOF {
					return
				}

				select {
				case entityRows <- types.EntityRow{Error: errdefs.WithStack(err)}:
					return
				case <-ctx.Done():
					return
				}
			}

			if len(line) < 2 {
				select {
				case entityRows <- types.EntityRow{Error: errdefs.Errorf("at least 2 values per row, got %d value(s) at row %d", len(line), i)}:
					return
				case <-ctx.Done():
					return
				}
			}

			unixMilli, err := strconv.Atoi(line[1])
			if err != nil {
				select {
				case entityRows <- types.EntityRow{Error: errdefs.WithStack(err)}:
					return
				case <-ctx.Done():
					return
				}
			}

			select {
			case entityRows <- types.EntityRow{
				EntityKey: line[0],
				UnixMilli: int64(unixMilli),
				Values:    line[2:],
				Error:     nil,
			}:
				// nothing to do
			case <-ctx.Done():
				return
			}
		}
	}()
	return entityRows, header, nil
}
