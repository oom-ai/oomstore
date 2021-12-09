package oomstore

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/ethhte88/oomstore/internal/database/metadata"
	"github.com/ethhte88/oomstore/internal/database/offline"
	"github.com/ethhte88/oomstore/pkg/oomstore/types"
)

// Import data into the offline feature store as a new revision.
// In the future we want to support more diverse data sources.
func (s *OomStore) Import(ctx context.Context, opt types.ImportOpt) (int, error) {
	importOpt, err := s.parseImportOpt(ctx, opt)
	if err != nil {
		return 0, err
	}
	switch dataSource := opt.DataSource.(type) {
	case types.CsvFileDataSource:
		file, err := os.Open(dataSource.InputFilePath)
		if err != nil {
			return 0, err
		}
		defer file.Close()
		importOpt.dataSource = types.CsvReaderDataSource{
			Reader:    file,
			Delimiter: dataSource.Delimiter,
		}
		return s.csvReaderImport(ctx, importOpt)
	case types.CsvReaderDataSource:
		return s.csvReaderImport(ctx, importOpt)
	case types.TableLinkDataSource:
		return s.tableLinkImport(ctx, importOpt)
	default:
		return 0, fmt.Errorf("unsupported data source: %T", opt.DataSource)
	}
}

func (s *OomStore) csvReaderImport(ctx context.Context, opt *importOpt) (int, error) {
	dataSource := opt.dataSource.(types.CsvReaderDataSource)
	// make sure csv data source has all defined columns
	csvReader := csv.NewReader(dataSource.Reader)
	csvReader.Comma = []rune(dataSource.Delimiter)[0]

	header, err := csvReader.Read()
	if err != nil {
		return 0, err
	}
	if hasDup(header) {
		return 0, fmt.Errorf("csv data source has duplicated columns: %v", header)
	}
	columnNames := append([]string{opt.entity.Name}, opt.features.Names()...)
	if !stringSliceEqual(header, columnNames) {
		return 0, fmt.Errorf("csv header of the data source %v doesn't match the feature group schema %v", header, columnNames)
	}

	newRevisionID, dataTableName, err := s.metadata.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:    0,
		GroupID:     opt.group.ID,
		DataTable:   nil,
		Description: opt.description,
		Anchored:    opt.revision != nil,
	})
	if err != nil {
		return 0, err
	}

	revision, err := s.offline.Import(ctx, offline.ImportOpt{
		Entity:        opt.entity,
		Features:      opt.features,
		Header:        header,
		Revision:      opt.revision,
		CsvReader:     csvReader,
		DataTableName: dataTableName,
	})
	if err != nil {
		return 0, err
	}

	if opt.revision != nil {
		revision = *opt.revision
	}
	if err := s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
		RevisionID:  newRevisionID,
		NewRevision: &revision,
	}); err != nil {
		return 0, err
	}

	// TODO: clean up revision and data_table if import failed

	return newRevisionID, nil
}

func (s *OomStore) tableLinkImport(ctx context.Context, opt *importOpt) (int, error) {
	dataSource := opt.dataSource.(types.TableLinkDataSource)

	// Make sure all features existing with correct value type
	tableSchema, err := s.offline.TableSchema(ctx, dataSource.TableName)
	if err != nil {
		return 0, err
	}
	validate := func(f *types.Feature) error {
		for _, field := range tableSchema.Fields {
			if field.Name == f.Name {
				if field.ValueType != f.ValueType {
					return fmt.Errorf("expect value type '%s', got '%s'", f.ValueType, field.ValueType)
				}
				return nil
			}
		}
		return fmt.Errorf("field '%s' found in target table", f.Name)
	}
	for _, feature := range opt.features {
		if err := validate(feature); err != nil {
			return 0, err
		}
	}

	var revision int64
	if opt.revision == nil {
		revision = time.Now().UnixMilli()
	} else {
		revision = *opt.revision
	}
	newRevisionID, _, err := s.metadata.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:    revision,
		GroupID:     opt.group.ID,
		DataTable:   &dataSource.TableName,
		Description: opt.description,
		Anchored:    opt.revision != nil,
	})
	if err != nil {
		return 0, err
	}

	return newRevisionID, nil
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
	group, err := s.metadata.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return nil, err
	}

	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &group.ID,
	})
	if err != nil {
		return nil, err
	}
	if features == nil {
		err = fmt.Errorf("no features under group: %s", opt.GroupName)
		return nil, err
	}

	entity := group.Entity
	if entity == nil {
		return nil, fmt.Errorf("no entity found by group: %s", opt.GroupName)
	}

	return &importOpt{
		dataSource:  opt.DataSource,
		entity:      entity,
		group:       group,
		features:    features,
		revision:    opt.Revision,
		description: opt.Description,
	}, nil
}

type importOpt struct {
	dataSource  interface{}
	entity      *types.Entity
	group       *types.Group
	features    types.FeatureList
	revision    *int64
	description string
}
