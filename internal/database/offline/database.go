package offline

import (
	"context"
	"io"

	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetPointInTimeFeatureValues(ctx context.Context, opt GetPointInTimeFeatureValuesOpt) (map[string]dbutil.RowMap, error)
	GetFeatureValuesStream(ctx context.Context, opt GetFeatureValuesStreamOpt) (<-chan *types.RawFeatureValueRecord, error)
	ImportBatchFeatures(ctx context.Context, opt ImportBatchFeaturesOpt) (int64, string, error)

	ValueTypeTag(dbValueType string) (string, error)
	io.Closer
}
