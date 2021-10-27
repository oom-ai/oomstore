package offline

import (
	"context"
	"io"

	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetPointInTimeFeatureValues(ctx context.Context, entity *types.Entity, entityRows []types.EntityRow, revisionRanges []*types.RevisionRange, features []*types.RichFeature) (dataMap map[string]dbutil.RowMap, err error)
	GetFeatureValuesStream(ctx context.Context, opt GetFeatureValuesStreamOpt) (<-chan *types.RawFeatureValueRecord, error)
	ImportBatchFeatures(ctx context.Context, opt types.ImportBatchFeaturesOpt, entity *types.Entity, features []*types.Feature, header []string) (int64, string, error)

	ValueTypeTag(dbValueType string) (string, error)
	io.Closer
}
