package online

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
)

type Store interface {
	GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, featureNames []string) (database.RowMap, error)
	GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, entityKeys, featureNames []string) (map[string]database.RowMap, error)
}
