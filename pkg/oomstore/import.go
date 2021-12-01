package oomstore

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

// ChannelImport a CSV data source into the feature store as a new revision.
// In the future we want to support more diverse data sources.
func (s *OomStore) ChannelImport(ctx context.Context, opt types.ChannelImport) (int, error) {
	group, err := s.metadata.GetGroupByName(ctx, opt.GroupName)
	if err != nil {
		return 0, err
	}

	features, err := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{
		GroupID: &group.ID,
	})
	if err != nil {
		return 0, err
	}
	if features == nil {
		return 0, fmt.Errorf("no features under group: %s", opt.GroupName)
	}

	entity := group.Entity
	if entity == nil {
		return 0, fmt.Errorf("no entity found by group: %s", opt.GroupName)
	}

	// make sure csv data source has all defined columns
	csvReader := csv.NewReader(opt.DataSource.Reader)
	csvReader.Comma = []rune(opt.DataSource.Delimiter)[0]

	header, err := csvReader.Read()
	if err != nil {
		return 0, err
	}
	if hasDup(header) {
		return 0, fmt.Errorf("csv data source has duplicated columns: %v", header)
	}
	columnNames := append([]string{entity.Name}, features.Names()...)
	if !stringSliceEqual(header, columnNames) {
		return 0, fmt.Errorf("csv header of the data source %v doesn't match the feature group schema %v", header, columnNames)
	}

	var revision int64
	if opt.Revision != nil {
		revision = *opt.Revision
	}

	newRevisionID, dataTableName, err := s.metadata.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision: revision,
		GroupID:  group.ID,
		// TODO: support user-defined DataTable
		DataTable:   nil,
		Description: opt.Description,
		Anchored:    opt.Revision != nil,
	})
	if err != nil {
		return 0, err
	}

	revision, err = s.offline.Import(ctx, offline.ImportOpt{
		Entity:        entity,
		Features:      features,
		Header:        header,
		Revision:      opt.Revision,
		CsvReader:     csvReader,
		DataTableName: dataTableName,
	})
	if err != nil {
		return 0, err
	}

	if opt.Revision == nil {
		if err := s.metadata.UpdateRevision(ctx, metadata.UpdateRevisionOpt{
			RevisionID:  newRevisionID,
			NewRevision: &revision,
		}); err != nil {
			return 0, err
		}
	}

	// TODO: clean up revision and data_table if import failed

	return newRevisionID, nil
}

// Import is similar to ChannelImport, the only difference is that it takes in InputFilePath
// as an argument.
func (s *OomStore) Import(ctx context.Context, opt types.ImportOpt) (int, error) {
	file, err := os.Open(opt.DataSource.InputFilePath)
	if err != nil {
		return 0, err
	}
	defer file.Close()
	return s.ChannelImport(ctx, types.ChannelImport{
		GroupName:   opt.GroupName,
		Description: opt.Description,
		DataSource: types.CsvDataSource{
			Reader:    file,
			Delimiter: opt.DataSource.Delimiter,
		},
		Revision: opt.Revision,
	})
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
