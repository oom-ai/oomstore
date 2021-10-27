package cmd

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"strconv"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

type JoinHistoricalFeaturesOpt struct {
	InputFilePath string
	FeatureNames  []string
}

func joinHistoricalFeatures(ctx context.Context, store *oomstore.OomStore, opt JoinHistoricalFeaturesOpt) error {
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

	if err := outputJoinedHistoricalFeatures(featureRows, opt.FeatureNames); err != nil {
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

func outputJoinedHistoricalFeatures(featureRows []*types.EntityRowWithFeatures, featureNames []string) error {
	w := csv.NewWriter(os.Stdout)
	if err := w.Write(append([]string{"entity_key", "unix_time"}, featureNames...)); err != nil {
		return err
	}

	for _, row := range featureRows {
		record := []string{row.EntityKey, strconv.Itoa(int(row.UnixTime))}
		for _, featureKV := range row.FeatureValues {
			record = append(record, cast.ToString(featureKV.FeatureValue))
		}
		if err := w.Write(record); err != nil {
			return err
		}
	}
	w.Flush()
	return nil
}
