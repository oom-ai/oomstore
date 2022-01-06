package cmd

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/pkg/errors"
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
	entityRows, header, err := oomstore.GetEntityRowsFromInputFile(opt.InputFilePath)
	if err != nil {
		return err
	}

	joinResult, err := store.ChannelJoin(ctx, types.ChannelJoinOpt{
		FeatureFullNames: opt.FeatureNames,
		EntityRows:       entityRows,
		ValueNames:       header[2:],
	})
	if err != nil {
		return err
	}

	if err := printJoinResult(joinResult, output); err != nil {
		return err
	}

	return nil
}

func printJoinResult(joinResult *types.JoinResult, output string) error {
	switch output {
	case CSV:
		return printJoinResultInCSV(joinResult)
	case ASCIITable:
		return printJoinResultInASCIITable(joinResult)
	default:
		return errors.Errorf("unsupported output format %s", output)
	}
}

func printJoinResultInCSV(joinResult *types.JoinResult) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	if err := w.Write(joinResult.Header); err != nil {
		return err
	}
	for row := range joinResult.Data {
		if err := w.Write(joinRecord(row)); err != nil {
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
		table.Append(joinRecord(row))
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
