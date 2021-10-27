package onestore

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/onestore-ai/onestore/internal/database/metadata"
	"github.com/onestore-ai/onestore/internal/database/offline"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func getCsvHeader(filePath string) ([]string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	r := csv.NewReader(f)
	header, err := r.Read()
	if err != nil {
		return nil, err
	}
	return header, nil
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

func (s *OneStore) ImportBatchFeatures(ctx context.Context, opt types.ImportBatchFeaturesOpt) error {
	// get columns of the group
	columns, err := s.metadata.ListFeature(ctx, &opt.GroupName)
	if err != nil {
		return err
	}

	// get entity info
	group, err := s.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return err
	}
	entity, err := s.GetEntity(ctx, group.EntityName)
	if err != nil {
		return err
	}

	// make sure csv data source has all defined columns
	header, err := getCsvHeader(opt.DataSource.FilePath)
	if err != nil {
		return err
	}
	if hasDup(header) {
		return fmt.Errorf("csv data source has duplicated columns: %v", header)
	}
	columnNames := []string{entity.Name}
	for _, column := range columns {
		columnNames = append(columnNames, column.Name)
	}
	if !stringSliceEqual(header, columnNames) {
		return fmt.Errorf("csv header of the data source %v doesn't match the feature group schema %v", header, columnNames)
	}

	revision, dataTable, err := s.offline.Import(ctx, offline.ImportOpt{
		ImportBatchFeaturesOpt: opt,
		Entity:                 entity,
		Features:               columns,
		Header:                 header,
	})
	if err != nil {
		return err
	}

	return s.metadata.InsertRevision(ctx, metadata.InsertRevisionOpt{
		Revision:    revision,
		GroupName:   opt.GroupName,
		DataTable:   dataTable,
		Description: opt.Description,
	})
}
