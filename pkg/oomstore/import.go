package oomstore

import (
	"bufio"
	"context"
	"os"
	"time"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Import data into the offline feature store as a new revision.
// In the future we want to support more diverse data sources.
func (s *OomStore) Import(ctx context.Context, opt types.ImportOpt) (int, error) {
	importOpt, err := s.parseImportOpt(ctx, opt)
	if err != nil {
		return 0, err
	}
	switch opt.DataSourceType {
	case types.CSV_FILE:
		src := importOpt.CsvFileDataSource
		file, err := os.Open(src.InputFilePath)
		if err != nil {
			return 0, err
		}
		defer file.Close()
		return s.csvReaderImport(ctx, importOpt, &types.CsvReaderDataSource{
			Reader:    file,
			Delimiter: src.Delimiter,
		})
	case types.CSV_READER:
		return s.csvReaderImport(ctx, importOpt, opt.CsvReaderDataSource)
	case types.TABLE_LINK:
		return s.tableLinkImport(ctx, importOpt, opt.TableLinkDataSource)
	default:
		return 0, errdefs.Errorf("unsupported data source: %v", opt.DataSourceType)
	}
}

func (s *OomStore) csvReaderImport(ctx context.Context, opt *importOpt, dataSource *types.CsvReaderDataSource) (int, error) {
	//make sure csv data source has all defined columns
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

	newRevisionID, _, err := s.createRevision(ctx, metadata.CreateRevisionOpt{
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

func (s *OomStore) tableLinkImport(ctx context.Context, opt *importOpt, dataSource *types.TableLinkDataSource) (int, error) {
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
	newRevisionID, _, err := s.createRevision(ctx, metadata.CreateRevisionOpt{
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

func validateTableSchema(schema *types.DataTableSchema, features types.FeatureList) error {
	validate := func(f *types.Feature) error {
		for _, field := range schema.Fields {
			if field.Name == f.Name {
				if field.ValueType != f.ValueType {
					return errdefs.Errorf("expect value type '%s', got '%s'", f.ValueType, field.ValueType)
				}
				return nil
			}
		}
		return errdefs.Errorf("field '%s' found in target table", f.Name)
	}
	for _, feature := range features {
		if err := validate(feature); err != nil {
			return err
		}
	}
	return nil
}
func hasDup(a []string) bool {
	s := make(map[string]bool)
	for _, e := range a {
		if s[e] {
			return true
		}
		s[e] = true
	}
	return false
}

func readHeader(source *offline.CSVSource, expectedColumns []string) ([]string, error) {
	// read header does not need pass down features
	header, err := dbutil.ReadLine(dbutil.ReadLineOpt{
		Source: source,
	})
	if err != nil {
		return nil, err
	}
	if hasDup(cast.ToStringSlice(header)) {
		return nil, errdefs.Errorf("csv data source has duplicated columns: %v", header)
	}
	if !stringSliceEqual(cast.ToStringSlice(header), expectedColumns) {
		return nil, errdefs.Errorf("csv header of the data source %v doesn't match the feature group schema %v", header, expectedColumns)
	}
	return cast.ToStringSlice(header), nil
}

func stringSliceEqual(a, b []string) bool {
	ma := make(map[string]bool)
	mb := make(map[string]bool)
	for _, e := range a {
		ma[e] = true
	}
	for _, e := range b {
		mb[e] = true
	}
	if len(ma) != len(mb) {
		return false
	}
	for k := range mb {
		if _, ok := ma[k]; !ok {
			return false
		}
	}
	return true
}

func (s *OomStore) parseImportOpt(ctx context.Context, opt types.ImportOpt) (*importOpt, error) {
	entity, group, features, err := s.getGroupInfo(ctx, opt.GroupName)
	if err != nil {
		return nil, err
	}
	return &importOpt{
		ImportOpt:  &opt,
		entityName: entity.Name,
		group:      group,
		features:   features,
	}, nil
}

func (s *OomStore) getGroupInfo(ctx context.Context, groupName string) (*types.Entity, *types.Group, types.FeatureList, error) {
	group, err := s.metadata.GetGroupByName(ctx, groupName)
	if err != nil {
		return nil, nil, nil, err
	}

	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &group.ID,
	})
	if err != nil {
		return nil, nil, nil, err
	}
	if features == nil {
		err = errdefs.Errorf("no features under group: %s", groupName)
		return nil, nil, nil, err
	}

	entity := group.Entity
	if entity == nil {
		return nil, nil, nil, errdefs.Errorf("no entity found by group: %s", groupName)
	}
	return entity, group, features, nil
}

type importOpt struct {
	*types.ImportOpt
	entityName string
	group      *types.Group
	features   types.FeatureList
}
