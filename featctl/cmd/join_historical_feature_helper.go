package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"io"
	"os"
	"strconv"

	"github.com/olekukonko/tablewriter"
	"github.com/oom-ai/oomstore/internal/database/metadata"
	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

type JoinHistoricalFeaturesOpt struct {
	InputFilePath string
	FeatureNames  []string
}

func joinHistoricalFeatures(ctx context.Context, store *oomstore.OomStore, opt JoinHistoricalFeaturesOpt, output string) error {
	entityRows, err := getEntityRowsFromInputFile(opt.InputFilePath)
	if err != nil {
		return err
	}

	features := store.ListFeature(ctx, metadata.ListFeatureOpt{
		FeatureNames: &opt.FeatureNames,
	})
	if err != nil {
		return nil
	}

	joinResult, err := store.GetHistoricalFeatureValues(ctx, types.GetHistoricalFeatureValuesOpt{
		FeatureIDs: features.Ids(),
		EntityRows: entityRows,
	})
	if err != nil {
		return err
	}

	if err := printJoinedHistoricalFeatures(joinResult, output); err != nil {
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

func printJoinedHistoricalFeatures(joinResult *types.JoinResult, output string) error {
	switch output {
	case CSV:
		return printJoinedHistoricalFeaturesInCSV(joinResult)
	case ASCIITable:
		return printJoinedHistoricalFeaturesInASCIITable(joinResult)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printJoinedHistoricalFeaturesInCSV(joinResult *types.JoinResult) error {
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

func printJoinedHistoricalFeaturesInASCIITable(joinResult *types.JoinResult) error {
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
