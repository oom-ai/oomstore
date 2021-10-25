package offline

import (
	"context"
	"fmt"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/offline/postgres"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetPointInTimeFeatureValues(ctx context.Context, entity *types.Entity, revisionRanges []*types.RevisionRange, features []*types.RichFeature, entityRows []types.EntityRow) (dataMap map[string]database.RowMap, err error)
	LoadLocalFile(ctx context.Context, filePath, tableName, delimiter string, header []string) error
	GetFeatureValuesStream(ctx context.Context, opt types.GetFeatureValuesStreamOpt) (<-chan *types.RawFeatureValueRecord, error)
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
