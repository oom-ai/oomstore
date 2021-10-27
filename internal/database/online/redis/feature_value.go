package redis

import (
	"context"

	"github.com/onestore-ai/onestore/internal/database/dbutil"
	"github.com/onestore-ai/onestore/internal/database/online"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
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

	rowMap := make(dbutil.RowMap)
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
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	res := make(map[string]dbutil.RowMap)
	for _, entityKey := range opt.EntityKeys {
		rowMap, err := db.Get(ctx, online.GetOpt{
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
