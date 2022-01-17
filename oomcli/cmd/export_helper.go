package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func export(ctx context.Context, store *oomstore.OomStore, opt types.ChannelExportOpt, output string) error {
	exportResult, err := store.ChannelExport(ctx, opt)
	if err != nil {
		return err
	}
	if err := printExportResult(exportResult, output); err != nil {
		return fmt.Errorf("failed printing historical features: %+v", err)
	}
	return nil
}

func printExportResult(exportResult *types.ExportResult, output string) error {
	switch output {
	case CSV:
		return printExportResultInCSV(exportResult)
	case ASCIITable:
		return printExportResultInASCIITable(exportResult)
	default:
		return errdefs.Errorf("unsupported output format %s", output)
	}
}

func printExportResultInCSV(exportResult *types.ExportResult) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()

	if err := w.Write(exportResult.Header); err != nil {
		return err
	}

	for row := range exportResult.Data {
		if err := w.Write(cast.ToStringSlice([]interface{}(row))); err != nil {
			return err
		}
	}
	return exportResult.CheckStreamError()
}

func printExportResultInASCIITable(exportResult *types.ExportResult) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(exportResult.Header)
	table.SetAutoFormatHeaders(false)

	for row := range exportResult.Data {
		table.Append(cast.ToStringSlice([]interface{}(row)))
	}

	table.Render()
	return exportResult.CheckStreamError()
}
