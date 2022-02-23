package oomstore

import (
	"bufio"
	"context"
	"time"

	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// csvReaderImportBatch imports batch feature from external data source to offline store through csv reader.
func (s *OomStore) csvReaderImportBatch(ctx context.Context, opt *importOpt, dataSource *types.CsvReaderDataSource) (int, error) {
	// make sure csv data source has all defined columns
	reader := bufio.NewReader(dataSource.Reader)
	source := &offline.CSVSource{
		Reader:    reader,
		Delimiter: dataSource.Delimiter,
	}
	// read header
	columnNames := append([]string{opt.entityName}, opt.features.Names()...)
	header, err := readHeader(source, columnNames)
	if err != nil {
		return 0, err
	}

	newRevisionID, err := s.createRevision(ctx, metadata.CreateRevisionOpt{
		Revision:      time.Now().UnixMilli(),
		GroupID:       opt.group.ID,
		SnapshotTable: nil,
		Description:   opt.Description,
		Anchored:      opt.Revision != nil,
	})
	if err != nil {
		return 0, err
	}
	snapshotTableName := dbutil.OfflineBatchSnapshotTableName(opt.group.ID, int64(newRevisionID))
	revision, err := s.offline.Import(ctx, offline.ImportOpt{
		EntityName:        opt.entityName,
		Features:          opt.features,
		Header:            cast.ToStringSlice(header),
		Revision:          opt.Revision,
		SnapshotTableName: snapshotTableName,
		Source:            source,
		Category:          types.CategoryBatch,
	})
	if err != nil {
		return 0, err
	}

	if opt.Revision != nil {
		revision = *opt.Revision
	}
	if err = s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
		RevisionID:       newRevisionID,
		NewRevision:      &revision,
		NewSnapshotTable: &snapshotTableName,
	}); err != nil {
		return 0, err
	}

	// TODO: clean up revision and data_table if import failed

	return newRevisionID, nil
}

// tableLinkImportBatch links external data source to offline store for batch feature.
func (s *OomStore) tableLinkImportBatch(ctx context.Context, opt *importOpt, dataSource *types.TableLinkDataSource) (int, error) {
	// Make sure all features existing with correct value type
	tableSchema, err := s.offline.TableSchema(ctx, offline.TableSchemaOpt{
		TableName: dataSource.TableName,
	})
	if err != nil {
		return 0, err
	}
	if err = validateTableSchema(tableSchema, opt.features); err != nil {
		return 0, err
	}
	var revision int64
	if opt.Revision == nil {
		revision = time.Now().UnixMilli()
	} else {
		revision = *opt.Revision
	}
	newRevisionID, err := s.createRevision(ctx, metadata.CreateRevisionOpt{
		Revision:      revision,
		GroupID:       opt.group.ID,
		SnapshotTable: &dataSource.TableName,
		Description:   opt.Description,
		Anchored:      opt.Revision != nil,
	})
	if err != nil {
		return 0, err
	}

	return newRevisionID, nil
}
