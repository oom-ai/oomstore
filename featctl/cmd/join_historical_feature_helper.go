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

type JoinHistoricalFeaturesOpt struct {
	InputFilePath string
	FeatureNames  []string
}

func joinHistoricalFeatures(ctx context.Context, store *oomstore.OomStore, opt JoinHistoricalFeaturesOpt, output string) error {
	entityRows, err := getEntityRowsFromInputFile(opt.InputFilePath)
	if err != nil {
		return err
	}

	featureRows, err := store.GetHistoricalFeatureValues(ctx, types.GetHistoricalFeatureValuesOpt{
		FeatureNames: opt.FeatureNames,
		EntityRows:   entityRows,
	})
	if err != nil {
		return err
	}

	if err := printJoinedHistoricalFeatures(featureRows, opt.FeatureNames, output); err != nil {
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

func printJoinedHistoricalFeatures(featureRows <-chan *types.EntityRowWithFeatures, featureNames []string, output string) error {
	switch output {
	case CSV:
		return printJoinedHistoricalFeaturesInCSV(featureRows, featureNames)
	case ASCIITable:
		return printJoinedHistoricalFeaturesInASCIITable(featureRows, featureNames)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printJoinedHistoricalFeaturesInCSV(featureRows <-chan *types.EntityRowWithFeatures, featureNames []string) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	if err := w.Write(joinHeader(featureNames)); err != nil {
		return err
	}
	for row := range featureRows {
		if err := w.Write(joinRecord(row)); err != nil {
			return err
		}
	}
	return nil
}

func printJoinedHistoricalFeaturesInASCIITable(featureRows <-chan *types.EntityRowWithFeatures, featureNames []string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(joinHeader(featureNames))
	table.SetAutoFormatHeaders(false)

	for row := range featureRows {
		table.Append(joinRecord(row))
	}
	table.Render()
	return nil
}

func joinHeader(featureNames []string) []string {
	return append([]string{"entity_key", "unix_time"}, featureNames...)
}

func joinRecord(row *types.EntityRowWithFeatures) []string {
	record := []string{row.EntityKey, strconv.Itoa(int(row.UnixTime))}
	for _, featureKV := range row.FeatureValues {
		record = append(record, cast.ToString(featureKV.FeatureValue))
	}
	return record
}
