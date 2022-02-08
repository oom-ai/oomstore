package oomstore

import (
	"context"
	"os"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// Import API imports data from external data source to offline store.
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
		source := &types.CsvReaderDataSource{
			Reader:    file,
			Delimiter: src.Delimiter,
		}
		if importOpt.group.Category == types.CategoryStream {
			return 0, s.csvReaderImportStream(ctx, importOpt, source)
		} else {
			return s.csvReaderImportBatch(ctx, importOpt, source)
		}
	case types.CSV_READER:
		if importOpt.group.Category == types.CategoryStream {
			return 0, s.csvReaderImportStream(ctx, importOpt, opt.CsvReaderDataSource)
		} else {
			return s.csvReaderImportBatch(ctx, importOpt, opt.CsvReaderDataSource)
		}
	case types.TABLE_LINK:
		if importOpt.group.Category == types.CategoryStream {
			return 0, s.tableLinkImportStream(ctx, importOpt, opt.TableLinkDataSource)
		} else {
			return s.tableLinkImportBatch(ctx, importOpt, opt.TableLinkDataSource)
		}
	default:
		return 0, errdefs.Errorf("unsupported data source: %v", opt.DataSourceType)
	}
}

type importOpt struct {
	*types.ImportOpt
	entityName string
	group      *types.Group
	features   types.FeatureList
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
		err = errdefs.Errorf("no features under group: %s", opt.GroupName)
		return nil, err
	}

	entity := group.Entity
	if entity == nil {
		return nil, errdefs.Errorf("no entity found by group: %s", opt.GroupName)
	}
	return &importOpt{
		ImportOpt:  &opt,
		entityName: entity.Name,
		group:      group,
		features:   features,
	}, nil
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
