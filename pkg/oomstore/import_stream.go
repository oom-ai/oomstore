package oomstore

import (
	"bufio"
	"context"
	"io"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore/util"

	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	ImportStreamBatchSize = 100
)

func (s *OomStore) ImportStream(ctx context.Context, opt types.ImportStreamOpt) error {
	switch opt.DataSourceType {
	case types.CSV_FILE:
		src := opt.CsvFileDataSource
		file, err := os.Open(src.InputFilePath)
		if err != nil {
			return err
		}
		defer file.Close()
		return s.csvReaderImportStream(ctx, opt, &types.CsvReaderDataSource{
			Reader:    file,
			Delimiter: src.Delimiter,
		})
	case types.CSV_READER:
		return s.csvReaderImportStream(ctx, opt, opt.CsvReaderDataSource)
	default:
		return errdefs.Errorf("unsupported data source: %v", opt.DataSourceType)
	}
}

func (s *OomStore) csvReaderImportStream(ctx context.Context, opt types.ImportStreamOpt, dataSource *types.CsvReaderDataSource) error {
	// get group information
	entity, group, features, err := s.getGroupInfo(ctx, opt.GroupName)
	if err != nil {
		return err
	}
	// read header
	reader := bufio.NewReader(dataSource.Reader)
	source := &offline.CSVSource{
		Reader:    reader,
		Delimiter: dataSource.Delimiter,
	}
	columnNames := append([]string{entity.Name, "unix_milli"}, features.Names()...)
	header, err := readHeader(source, columnNames)
	if err != nil {
		return err
	}

	// load data
	records := make([]types.StreamRecord, 0, ImportStreamBatchSize)
	for {
		line, err := dbutil.ReadLine(dbutil.ReadLineOpt{
			Source:   source,
			Entity:   entity,
			Header:   header,
			Features: features,
		})
		if errdefs.Cause(err) == io.EOF {
			break
		}
		if err != nil {
			return err
		}
		if len(line) != len(header) {
			continue
		}
		records = append(records, generateStreamRecord(line, header, group, features))
		if len(records) == ImportStreamBatchSize {
			if err := s.pushStreamingRecords(ctx, records, entity.Name, group, features); err != nil {
				return err
			}
			records = make([]types.StreamRecord, 0, ImportStreamBatchSize)
		}
	}
	return s.pushStreamingRecords(ctx, records, entity.Name, group, features)
}

func (s *OomStore) pushStreamingRecords(ctx context.Context, records []types.StreamRecord, entityName string, group *types.Group, features types.FeatureList) error {
	buckets := make(map[int64][]types.StreamRecord)
	for _, record := range records {
		revision := lastRevisionForStream(int64(group.SnapshotInterval), record.UnixMilli)
		if _, ok := buckets[revision]; !ok {
			buckets[revision] = make([]types.StreamRecord, 0)
		}
		buckets[revision] = append(buckets[revision], record)
	}
	for revision, streamRecords := range buckets {
		pushOpt := offline.PushOpt{
			GroupID:      group.ID,
			Revision:     revision,
			EntityName:   entityName,
			FeatureNames: features.Names(),
			Records:      streamRecords,
		}

		if err := s.offline.Push(ctx, pushOpt); err != nil {
			if !errdefs.IsNotFound(err) {
				return err
			}

			if err = s.newRevisionForStream(ctx, group.ID, revision); err != nil {
				return err
			}
			// push data to new offline stream cdc table
			if err = s.offline.Push(ctx, pushOpt); err != nil {
				return err
			}
		}
	}
	return nil
}
func generateStreamRecord(line []interface{}, header []string, group *types.Group, features types.FeatureList) types.StreamRecord {
	var (
		entityKey string
		unixMilli int64
		values    []interface{}
	)
	for i, value := range line {
		if header[i] == group.Entity.Name {
			entityKey = cast.ToString(value)
		} else if header[i] == "unix_milli" {
			unixMilli = cast.ToInt64(value)
		}
	}
	for _, f := range features {
		idx := util.SliceIndex(len(header), func(i int) bool {
			return header[i] == f.Name
		})
		values = append(values, line[idx])
	}

	return types.StreamRecord{
		GroupID:   group.ID,
		EntityKey: entityKey,
		UnixMilli: unixMilli,
		Values:    values,
	}
}
