package offline

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetPointInTimeFeatureValues(ctx context.Context, features []*types.RichFeature, entityRows []types.EntityRow) (dataMap map[string]database.RowMap, err error)
}
