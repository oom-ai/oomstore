package redis

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
)

func (db *DB) GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, featureNames []string) (database.RowMap, error) {
	return make(database.RowMap), nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, entityKeys, featureNames []string) (map[string]database.RowMap, error) {
	return make(map[string]database.RowMap), nil
}
