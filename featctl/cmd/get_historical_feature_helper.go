package cmd

import (
	"context"
	"encoding/csv"
	"os"

	"github.com/onestore-ai/onestore/pkg/onestore"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
	"github.com/spf13/cast"
)

type getHistoricalFeatureOption struct {
	GroupName     string
	GroupRevision *int64
	FeatureNames  []string
	Limit         *uint64
}

func getHistoricalFeature(ctx context.Context, store *onestore.OneStore, opt getHistoricalFeatureOption) error {
	group, err := store.GetFeatureGroup(ctx, opt.GroupName)
	if err != nil {
		return err
	}

	if opt.GroupRevision != nil {
		group.Revision = opt.GroupRevision
	}
	walkOpt := types.WalkFeatureValuesOpt{
		FeatureGroup: *group,
		FeatureNames: opt.FeatureNames,
		Limit:        opt.Limit,
	}
	w := csv.NewWriter(os.Stdout)
	headerRow := true
	walkOpt.WalkFeatureValuesFunc = func(header []string, key string, values []interface{}) error {
		if headerRow {
			if err := w.Write(header); err != nil {
				return err
			}
			headerRow = false
		}
		record := []string{key}
		record = append(record, cast.ToStringSlice(values)...)
		return w.Write(record)
	}
	err = store.WalkFeatureValues(ctx, walkOpt)
	w.Flush()
	return err
}
