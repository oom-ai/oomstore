package online

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

type Store interface {
	GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, featureNames []string) (database.RowMap, error)
	GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, entityKeys, featureNames []string) (map[string]database.RowMap, error)
	LoadFeatureValuesStream(ctx context.Context, stream <-chan []interface{}, features []*types.Feature, revision *types.Revision) error
}
