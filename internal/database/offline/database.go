package offline

import (
	"context"
	"fmt"
	"io"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/offline/postgres"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetPointInTimeFeatureValues(ctx context.Context, entity *types.Entity, revisionRanges []*types.RevisionRange, features []*types.RichFeature, entityRows []types.EntityRow) (dataMap map[string]database.RowMap, err error)
	GetFeatureValuesStream(ctx context.Context, opt types.GetFeatureValuesStreamOpt) (<-chan *types.RawFeatureValueRecord, error)
	ImportBatchFeatures(ctx context.Context, opt types.ImportBatchFeaturesOpt, entity *types.Entity, features []*types.Feature, header []string) (int64, string, error)

	ValueTypeTag(dbValueType string) (string, error)
	io.Closer
}

var _ Store = &postgres.DB{}

func Open(opt types.OfflineStoreOpt) (Store, error) {
	switch opt.Backend {
	case types.POSTGRES:
		return postgres.Open(opt.PostgresDbOpt)
	default:
		return nil, fmt.Errorf("unsupported backend: %s", opt.Backend)
	}
}
