package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

type JoinOpt struct {
	InputFilePath string
	FeatureNames  []string
}

func join(ctx context.Context, store *oomstore.OomStore, opt JoinOpt, output string) error {
	entityRows, err := getEntityRowsFromInputFile(opt.InputFilePath)
	if err != nil {
		return err
	}

	features, err := store.ListFeature(ctx, types.ListFeatureOpt{
		FeatureNames: &opt.FeatureNames,
	})
	if err != nil {
		return nil
	}

	joinResult, err := store.Join(ctx, types.JoinOpt{
		FeatureIDs: features.IDs(),
		EntityRows: entityRows,
	})
	if err != nil {
		return err
	}

	if err := printJoinResult(joinResult, output); err != nil {
		return err
	}

	return nil
}

func getEntityRowsFromInputFile(inputFilePath string) (<-chan types.EntityRow, error) {
	input, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}
	entityRows := make(chan types.EntityRow)
	var readErr error
	go func() {
		defer close(entityRows)
		defer input.Close()
		reader := csv.NewReader(input)
		var i int64
		for {
			line, err := reader.Read()
			if err == io.EOF {
				break
			}
			if err != nil {
				readErr = err
				return
			}
			if len(line) != 2 {
				readErr = fmt.Errorf("expected 2 values per row, got %d value(s) at row %d", len(line), i)
				return
			}
			unixTime, err := strconv.Atoi(line[1])
			if err != nil {
				readErr = err
				return
			}
			entityRows <- types.EntityRow{
				EntityKey: line[0],
				UnixTime:  int64(unixTime),
			}
			i++
		}
	}()
	if readErr != nil {
		return nil, readErr
	}
	return entityRows, nil
}

func printJoinResult(joinResult *types.JoinResult, output string) error {
	switch output {
	case CSV:
		return printJoinResultInCSV(joinResult)
	case ASCIITable:
		return printJoinResultInASCIITable(joinResult)
	default:
		return fmt.Errorf("unsupported output format %s", output)
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
