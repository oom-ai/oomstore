package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

func getHistoricalFeature(ctx context.Context, store *oomstore.OomStore, opt types.ExportFeatureValuesOpt, output string) error {
	fields, stream, err := store.ExportFeatureValues(ctx, opt)
	if err != nil {
		return err
	}

	if err := printHistoricalFeatures(fields, stream, output); err != nil {
		return fmt.Errorf("failed printing historical features: %+v", err)
	}
	return nil
}

func printHistoricalFeatures(fields []string, stream <-chan *types.RawFeatureValueRecord, output string) error {
	switch output {
	case CSV:
		return printHistoricalFeaturesInCSV(fields, stream)
	case ASCIITable:
		return printHistoricalFeaturesInASCIITable(fields, stream)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printHistoricalFeaturesInCSV(fields []string, stream <-chan *types.RawFeatureValueRecord) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	if err := w.Write(fields); err != nil {
		return err
	}
	for item := range stream {
		if item.Error != nil {
			return item.Error
		}
		if err := w.Write(cast.ToStringSlice(item.Record)); err != nil {
			return err
		}
	}
	return nil
}

func printHistoricalFeaturesInASCIITable(fields []string, stream <-chan *types.RawFeatureValueRecord) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(fields)
	table.SetAutoFormatHeaders(false)

	for item := range stream {
		if item.Error != nil {
			return item.Error
		}
		table.Append(cast.ToStringSlice(item.Record))
	}

	table.Render()
	return nil
}
