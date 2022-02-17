package cmd

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/oom-ai/oomstore/pkg/errdefs"
	"github.com/spf13/cast"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

type JoinOpt struct {
	InputFilePath string
	FeatureNames  []string
}

func join(ctx context.Context, store *oomstore.OomStore, opt JoinOpt, output string) error {
	ctx, cancel := context.WithCancel(ctx)
	defer cancel()

	joinResult, err := store.Join(ctx, types.JoinOpt{
		FeatureNames:  opt.FeatureNames,
		InputFilePath: opt.InputFilePath,
	})
	if err != nil {
		return err
	}

	return printJoinResult(joinResult, output)
}

func printJoinResult(joinResult *types.JoinResult, output string) error {
	switch output {
	case CSV:
		return printJoinResultInCSV(joinResult)
	case ASCIITable:
		return printJoinResultInASCIITable(joinResult)
	default:
		return errdefs.Errorf("unsupported output format %s", output)
	}
}

func printJoinResultInCSV(joinResult *types.JoinResult) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	if err := w.Write(joinResult.Header); err != nil {
		return err
	}
	for row := range joinResult.Data {
		if row.Error != nil {
			return row.Error
		}
		if err := w.Write(joinRecord(row.Record)); err != nil {
			return err
		}
	}
	return nil
}

func printJoinResultInASCIITable(joinResult *types.JoinResult) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(joinResult.Header)
	table.SetAutoFormatHeaders(false)

	for row := range joinResult.Data {
		if row.Error != nil {
			return row.Error
		}
		table.Append(joinRecord(row.Record))
	}
	table.Render()
	return nil
}

func joinRecord(row []interface{}) []string {
	record := make([]string, 0, len(row))
	for _, value := range row {
		record = append(record, cast.ToString(value))
	}
	return record
}
