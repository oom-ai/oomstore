package redis

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
)

func (db *DB) GetFeatureValues(ctx context.Context, dataTable, entityName, entityKey string, revisionId int32, featureNames []string) (database.RowMap, error) {
	rowMap := make(database.RowMap)

	key, err := SerializeRedisKey(revisionId, entityKey)
	if err != nil {
		return rowMap, err
	}

	values, err := db.HMGet(ctx, key, featureNames...).Result()
	if err != nil {
		return rowMap, err
	}
	for i, v := range values {
		rowMap[featureNames[i]] = v
	}
	return rowMap, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) GetFeatureValuesWithMultiEntityKeys(ctx context.Context, dataTable, entityName string, revisionId int32, entityKeys, featureNames []string) (map[string]database.RowMap, error) {
	res := make(map[string]database.RowMap)
	for _, entityKey := range entityKeys {
		rowMap, err := db.GetFeatureValues(ctx, dataTable, entityName, entityKey, revisionId, featureNames)
		if err != nil {
			return res, err
		}
		res[entityKey] = rowMap
	}
	return res, nil
}
