package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
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

func getEntityRowsFromInputFile(inputFilePath string) ([]types.EntityRow, error) {
	input, err := os.Open(inputFilePath)
	if err != nil {
		return nil, err
	}
	defer input.Close()

	lines, err := csv.NewReader(input).ReadAll()
	if err != nil {
		return nil, err
	}
	entityRows := make([]types.EntityRow, 0, len(lines))
	for i, line := range lines {
		if len(line) != 2 {
			return nil, fmt.Errorf("expected 2 values per row, got %d value(s) at row %d", len(line), i)
		}
		unixTime, err := strconv.Atoi(line[1])
		if err != nil {
			return nil, err
		}
		entityRows = append(entityRows, types.EntityRow{
			EntityKey: line[0],
			UnixTime:  int64(unixTime),
		})
	}
	return entityRows, nil
}

func printJoinedHistoricalFeatures(featureRows []*types.EntityRowWithFeatures, featureNames []string, output string) error {
	switch output {
	case CSV:
		return printJoinedHistoricalFeaturesInCSV(featureRows, featureNames)
	case ASCIITable:
		return printJoinedHistoricalFeaturesInASCIITable(featureRows, featureNames)
	default:
		return fmt.Errorf("unsupported output format %s", output)
	}
}

func printJoinedHistoricalFeaturesInCSV(featureRows []*types.EntityRowWithFeatures, featureNames []string) error {
	w := csv.NewWriter(os.Stdout)
	defer w.Flush()
	if err := w.Write(joinHeader(featureNames)); err != nil {
		return err
	}
	for _, row := range featureRows {
		if err := w.Write(joinRecord(row)); err != nil {
			return err
		}
	}
	return nil
}

func printJoinedHistoricalFeaturesInASCIITable(featureRows []*types.EntityRowWithFeatures, featureNames []string) error {
	table := tablewriter.NewWriter(os.Stdout)
	table.SetHeader(joinHeader(featureNames))
	table.SetAutoFormatHeaders(false)

	for _, row := range featureRows {
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
