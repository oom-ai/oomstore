package redis

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database"
	dbtypes "github.com/onestore-ai/onestore/internal/database/types"
	"github.com/onestore-ai/onestore/pkg/onestore/types"
)

func (db *DB) Get(ctx context.Context, opt types.GetFeatureValuesOpt) (database.RowMap, error) {
	key, err := SerializeRedisKey(opt.RevisionId, opt.EntityKey)
	if err != nil {
		return nil, err
	}

	featureIds := []string{}
	for _, f := range opt.Features {
		id, err := SerializeByValue(f.ID)
		if err != nil {
			return nil, err
		}
		featureIds = append(featureIds, id)
	}

	values, err := db.HMGet(ctx, key, featureIds...).Result()
	if err != nil {
		return nil, err
	}

	rowMap := make(database.RowMap)
	for i, v := range values {
		typedValue, err := DeserializeByTag(v, opt.Features[i].ValueType)
		if err != nil {
			return nil, err
		}
		rowMap[opt.Features[i].Name] = typedValue
	}
	return rowMap, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGetOnlineFeatureValues(ctx context.Context, opt dbtypes.MultiGetOnlineFeatureValuesOpt) (map[string]database.RowMap, error) {
	res := make(map[string]database.RowMap)
	for _, entityKey := range opt.EntityKeys {
		rowMap, err := db.Get(ctx, types.GetFeatureValuesOpt{
			DataTable:  opt.DataTable,
			EntityName: opt.EntityName,
			RevisionId: opt.RevisionId,
			EntityKey:  entityKey,
			Features:   opt.Features,
		})
		if err != nil {
			return res, err
		}
		res[entityKey] = rowMap
	}
	return res, nil
}
