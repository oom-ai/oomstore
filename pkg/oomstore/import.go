package oomstore

import (
	"context"
	"encoding/csv"
	"fmt"

	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/internal/database/offline"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

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

func (s *OomStore) ImportBatchFeatures(ctx context.Context, opt types.ImportBatchFeaturesOpt) (int32, error) {
	// get columns of the group
	features := s.metadata.ListFeature(ctx, metadata.ListFeatureOpt{GroupID: &opt.GroupID})
	if features == nil {
		return 0, fmt.Errorf("no featues under group id: '%d'", opt.GroupID)
	}

	// get entity info
	group, err := s.GetFeatureGroup(ctx, opt.GroupID)
	if err != nil {
		return 0, err
	}
	entity := group.Entity
	if entity == nil {
		return 0, fmt.Errorf("no entity found by group id: '%d'", opt.GroupID)
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

	revision, dataTable, err := s.offline.Import(ctx, offline.ImportOpt{
		GroupName: group.Name,
		Entity:    entity,
		Features:  features,
		Header:    header,
		Revision:  opt.Revision,
		CsvReader: csvReader,
	})
	if err != nil {
		return 0, err
	}

	newRevisionID, err := s.metadata.CreateRevision(ctx, metadata.CreateRevisionOpt{
		Revision:    revision,
		GroupID:     group.ID,
		DataTable:   &dataTable,
		Description: opt.Description,
		Anchored:    opt.Revision != nil,
	})
	if err != nil {
		return 0, err
	}
	return newRevisionID, nil
}
