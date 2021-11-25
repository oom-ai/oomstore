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

func export(ctx context.Context, store *oomstore.OomStore, opt types.ExportOpt, output string) error {
	fields, stream, err := store.Export(ctx, opt)
	if err != nil {
		return err
	}

	if err := printExportResult(fields, stream, output); err != nil {
		return fmt.Errorf("failed printing historical features: %+v", err)
	}
	return nil
}

func printExportResult(fields []string, stream <-chan *types.ExportRecord, output string) error {
	switch output {
	case CSV:
		return printExportResultInCSV(fields, stream)
	case ASCIITable:
		return printExportResultInASCIITable(fields, stream)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printExportResultInCSV(fields []string, stream <-chan *types.ExportRecord) error {
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

func printExportResultInASCIITable(fields []string, stream <-chan *types.ExportRecord) error {
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
