package cmd

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/oom-ai/oomstore/pkg/oomstore"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
	"github.com/spf13/cast"
)

func getHistoricalFeature(ctx context.Context, store *oomstore.OomStore, opt types.ExportFeatureValuesOpt) error {
	fields, stream, err := store.ExportFeatureValues(ctx, opt)
	if err != nil {
		return err
	}

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
	return err
}
