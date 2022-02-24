package redis

import (
	"context"

	"github.com/oom-ai/oomstore/pkg/errdefs"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
	"github.com/oom-ai/oomstore/pkg/oomstore/types"
)

func (db *DB) Get(ctx context.Context, opt online.GetOpt) (dbutil.RowMap, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}

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
		id, err2 := dbutil.SerializeByValue(f.ID, Backend)
		if err2 != nil {
			return nil, err2
		}
		featureIDs = append(featureIDs, id)
	}

	values, err := db.HMGet(ctx, key, featureIDs...).Result()
	if err != nil {
		return nil, errdefs.WithStack(err)
	}

	rowMap := make(dbutil.RowMap)
	for i, v := range values {
		if v == nil {
			continue
		}
		deserializedValue, err := dbutil.DeserializeByValueType(v, opt.Features[i].ValueType, Backend)
		if err != nil {
			return nil, err
		}
		rowMap[opt.Features[i].FullName()] = deserializedValue
	}
	return rowMap, nil
}

// response: map[entity_key]map[feature_name]feature_value
func (db *DB) MultiGet(ctx context.Context, opt online.MultiGetOpt) (map[string]dbutil.RowMap, error) {
	if err := opt.Validate(); err != nil {
		return nil, err
	}

	res := make(map[string]dbutil.RowMap)
	for _, entityKey := range opt.EntityKeys {
		rowMap, err := db.Get(ctx, online.GetOpt{
			EntityKey:  entityKey,
			Group:      opt.Group,
			Features:   opt.Features,
			RevisionID: opt.RevisionID,
		})
		if err != nil {
			return res, errdefs.WithStack(err)
		}
		if len(rowMap) > 0 {
			res[entityKey] = rowMap
		}
	}
	return res, nil
}
