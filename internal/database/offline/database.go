package offline

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/internal/database/offline/postgres"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetPointInTimeFeatureValues(ctx context.Context, entity *types.Entity, revisionRanges []*types.RevisionRange, features []*types.RichFeature, entityRows []types.EntityRow) (dataMap map[string]database.RowMap, err error)
	GetFeatureValuesStream(ctx context.Context, opt types.GetFeatureValuesStreamOpt) (<-chan []interface{}, error)
}

var _ Store = &postgres.DB{}
