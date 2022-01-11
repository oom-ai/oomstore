package redis

import (
	"context"

	"github.com/pkg/errors"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	var (
		key string
		err error
	)

	if opt.Group.Category == types.CategoryBatch {
		key, err = serializeRedisKeyForBatchFeature(*opt.RevisionID, opt.EntityKey)
	} else {
		key, err = serializeRedisKeyForStreamFeature(opt.Group.ID, opt.EntityKey)
	}
	if err != nil {
		return nil, err
	}

	featureIDs := []string{}
	for _, f := range opt.Features {
		id, err := dbutil.SerializeByValue(f.ID, Backend)
		if err != nil {
			return nil, err
		}
		featureIDs = append(featureIDs, id)
	}

	values, err := db.HMGet(ctx, key, featureIDs...).Result()
	if err != nil {
		return nil, errors.WithStack(err)
	}

	rowMap := make(dbutil.RowMap)
	for i, v := range values {
		if v == nil {
			continue
		}
		typedValue, err := dbutil.DeserializeByValueType(v, opt.Features[i].ValueType, Backend)
		if err != nil {
			return nil, err
		}
		rowMap[opt.Features[i].FullName] = typedValue
	}
	return rowMap, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	res := make(map[string]dbutil.RowMap)
	for _, entityKey := range opt.EntityKeys {
		rowMap, err := db.Get(ctx, online.GetOpt{
			Entity:     opt.Entity,
			RevisionID: opt.RevisionID,
			Group:      opt.Group,
			EntityKey:  entityKey,
			Features:   opt.Features,
		})
		if err != nil {
			return res, errors.WithStack(err)
		}
		if len(rowMap) > 0 {
			res[entityKey] = rowMap
		}
	}
	return res, nil
}
