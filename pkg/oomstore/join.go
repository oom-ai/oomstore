package oomstore

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"math"
	"os"
	"sort"
	"strconv"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

// Get point-in-time correct feature values for each entity row.
// Currently, this API only supports batch features.
func (s *OomStore) ChannelJoin(ctx context.Context, opt types.ChannelJoinOpt) (*types.JoinResult, error) {
	if err := s.metadata.Refresh(); err != nil {
		return nil, fmt.Errorf("failed to refresh informer, err=%+v", err)
	}
	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureNames: &opt.FeatureNames,
	})

	features = features.Filter(func(f *types.Feature) bool {
		return f.Group.Category == types.BatchFeatureCategory
	})
	if len(features) == 0 {
		return nil, nil
	}

	entity, err := getSharedEntity(features)
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

// Get point-in-time correct feature values for each entity row.
// The method is similar to Join, except that both input and output are files on disk.
func (s *OomStore) Join(ctx context.Context, opt types.JoinOpt) error {
	entityRows, err := GetEntityRowsFromInputFile(opt.InputFilePath)
	if err != nil {
		return err
	}

	joinResult, err := s.ChannelJoin(ctx, types.ChannelJoinOpt{
		FeatureNames: opt.FeatureNames,
		EntityRows:   entityRows,
	})
	if err != nil {
		return err
	}
	return writeJoinResultToFile(opt.OutputFilePath, joinResult)
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
	revisions := s.metadata.ListRevision(ctx, &groupID)
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

func GetEntityRowsFromInputFile(inputFilePath string) (<-chan types.EntityRow, error) {
	input, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}
	entityRows := make(chan types.EntityRow)
	var readErr error
	go func() {
		defer close(entityRows)
		defer input.Close()
		reader := csv.NewReader(input)
		var i int64
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				readErr = err
				return
			}
			if len(line) != 2 {
				readErr = fmt.Errorf("expected 2 values per row, got %d value(s) at row %d", len(line), i)
				return
			}
			unixMilli, err := strconv.Atoi(line[1])
			if err != nil {
				readErr = err
				return
			}
			entityRows <- types.EntityRow{
				EntityKey: line[0],
				UnixMilli: int64(unixMilli),
			}
			i++
		}
	}()
	if readErr != nil {
		return nil, readErr
	}
	return entityRows, nil
}

func writeJoinResultToFile(outputFilePath string, joinResult *types.JoinResult) error {
	file, err := os.Create(outputFilePath)
	if err != nil {
		return err
	}
	defer file.Close()
	w := csv.NewWriter(file)
	defer w.Flush()

	if err := w.Write(joinResult.Header); err != nil {
		return err
	}
	for row := range joinResult.Data {
		if err := w.Write(joinRecord(row)); err != nil {
			return err
		}
	}
	return nil
}

func joinRecord(row []interface{}) []string {
	record := make([]string, 0, len(row))
	for _, value := range row {
		record = append(record, cast.ToString(value))
	}
	return record
}
