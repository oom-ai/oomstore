package export

import (
	"context"
	"encoding/csv"
	"fmt"
	"os"
	"time"

	"github.com/onestore-ai/onestore/pkg/onestore"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type ExportOpt struct {
	GroupName     string
	GroupRevision *int64
	FeatureNames  []string
	Limit         *uint64
}

func Export(ctx context.Context, store *onestore.OneStore, opt ExportOpt) error {
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
		for _, value := range values {
			if value == nil {
				record = append(record, "")
			} else if bytes, ok := value.([]byte); ok {
				record = append(record, string(bytes))
			} else if f, ok := value.(float64); ok {
				record = append(record, fmt.Sprintf("%f", f))
			} else if t, ok := value.(time.Time); ok {
				record = append(record, t.Format(time.RFC3339))
			}
		}
		return w.Write(record)
	}
	err = store.WalkFeatureValues(ctx, walkOpt)
	w.Flush()
	return err
}
