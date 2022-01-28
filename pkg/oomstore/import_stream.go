package oomstore

import (
	"bufio"
	"context"
	"io"

	"github.com/oom-ai/oomstore/pkg/oomstore/util"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

const (
	ImportStreamBatchSize = 100
)

func (s *OomStore) csvReaderImportStream(ctx context.Context, opt *importOpt, dataSource *types.CsvReaderDataSource) error {
	// read header
	reader := bufio.NewReader(dataSource.Reader)
	source := &offline.CSVSource{
		Reader:    reader,
		Delimiter: dataSource.Delimiter,
	}
	columnNames := append([]string{opt.entityName, "unix_milli"}, opt.features.Names()...)
	header, err := readHeader(source, columnNames)
	if err != nil {
		return err
	}

	// load data
	records := make([]types.StreamRecord, 0, ImportStreamBatchSize)
	for {
		line, err := dbutil.ReadLine(dbutil.ReadLineOpt{
			Source:     source,
			EntityName: opt.entityName,
			Header:     header,
			Features:   opt.features,
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
		records = append(records, generateStreamRecord(line, header, opt.group, opt.features))
		if len(records) == ImportStreamBatchSize {
			if err := s.pushStreamingRecords(ctx, records, opt.entityName, opt.group, opt.features); err != nil {
				return err
			}
			records = make([]types.StreamRecord, 0, ImportStreamBatchSize)
		}
	}
	return s.pushStreamingRecords(ctx, records, opt.entityName, opt.group, opt.features)
}

func (s *OomStore) tableLinkImportStream(ctx context.Context, opt *importOpt, dataSource *types.TableLinkDataSource) error {
	// get linked table schema
	tableSchema, err := s.offline.TableSchema(ctx, offline.TableSchemaOpt{
		TableName:      dataSource.TableName,
		CheckTimeRange: true,
	})
	if err != nil {
		return err
	}
	// validation
	if err = validateTableSchema(tableSchema, opt.features); err != nil {
		return err
	}
	if err = s.validateRevisions(ctx, opt.group.ID, tableSchema); err != nil {
		return err
	}

	_, err = s.createRevision(ctx, metadata.CreateRevisionOpt{
		Revision:    *tableSchema.TimeRange.MinUnixMilli,
		GroupID:     opt.group.ID,
		CdcTable:    &dataSource.TableName,
		Description: opt.Description,
		Anchored:    true,
	})
	if err != nil {
		return err
	}
	return nil
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

func (s *OomStore) validateRevisions(ctx context.Context, groupID int, schema *types.DataTableSchema) error {
	revisions, err := s.ListRevision(ctx, &groupID)
	if err != nil {
		return err
	}
	if schema.TimeRange.MinUnixMilli == nil || schema.TimeRange.MaxUnixMilli == nil {
		return errdefs.Errorf("table is empty")
	}
	revisionBeforeMin := revisions.Before(*schema.TimeRange.MinUnixMilli)
	revisionBeforeMax := revisions.Before(*schema.TimeRange.MaxUnixMilli)
	if revisionBeforeMax != revisionBeforeMin {
		return errdefs.Errorf("data table crosses with another offline table %s", revisionBeforeMax.CdcTable)
	}
	if revisionBeforeMin != nil {
		beforeTableSchema, err := s.offline.TableSchema(ctx, offline.TableSchemaOpt{
			TableName:      revisionBeforeMin.CdcTable,
			CheckTimeRange: true,
		})
		if err != nil {
			return err
		}
		if beforeTableSchema.TimeRange.MaxUnixMilli != nil && *beforeTableSchema.TimeRange.MaxUnixMilli >= *schema.TimeRange.MinUnixMilli {
			return errdefs.Errorf("data table crosses with another offline table %s", revisionBeforeMin.CdcTable)
		}
	}
	return nil
}
