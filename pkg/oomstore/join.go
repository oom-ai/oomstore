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

// ChannelJoin gets point-in-time correct feature values for each entity row.
// Currently, this API only supports batch features.
func (s *OomStore) ChannelJoin(ctx context.Context, opt types.ChannelJoinOpt) (*types.JoinResult, error) {
	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureFullNames: &opt.FeatureFullNames,
	})
	if err != nil {
		return nil, err
	}

	features = features.Filter(func(f *types.Feature) bool {
		return f.Group.Category == types.CategoryBatch
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
	revisionRangeMap := make(map[string][]*offline.RevisionRange)
	for groupName, featureList := range featureMap {
		if len(featureList) == 0 {
			continue
		}
		revisionRanges, err := s.buildRevisionRanges(ctx, featureList[0].Group)
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
		ValueNames:       opt.ValueNames,
	})
}

// Join gets point-in-time correct feature values for each entity row.
// The method is similar to Join, except that both input and output are files on disk.
// Input File should contain header, the first two columns of Input File should be
// entity_key, unix_milli, then followed by other real-time feature values.
func (s *OomStore) Join(ctx context.Context, opt types.JoinOpt) error {
	entityRows, header, err := GetEntityRowsFromInputFile(opt.InputFilePath)
	if err != nil {
		return err
	}

	joinResult, err := s.ChannelJoin(ctx, types.ChannelJoinOpt{
		FeatureFullNames: opt.FeatureFullNames,
		EntityRows:       entityRows,
		ValueNames:       header[2:],
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

func GetEntityRowsFromInputFile(inputFilePath string) (<-chan types.EntityRow, []string, error) {
	input, err := os.Open(inputFilePath)
	if err != nil {
		return nil, nil, err
	}
	reader := csv.NewReader(input)
	header, err := reader.Read()
	if err != nil {
		return nil, nil, err
	}
	entityRows := make(chan types.EntityRow)
	var readErr error
	i := 1
	go func() {
		defer close(entityRows)
		defer input.Close()
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				readErr = err
				return
			}
			if len(line) < 2 {
				readErr = fmt.Errorf("at least 2 values per row, got %d value(s) at row %d", len(line), i)
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
				Values:    line[2:],
			}
			i++
		}
	}()
	if readErr != nil {
		return nil, nil, readErr
	}
	return entityRows, header, nil
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
