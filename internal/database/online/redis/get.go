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

	key, err := serializeRedisKey(opt.Group, opt.EntityKey, opt.RevisionID)
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

func (db *DB) GetByGroup(ctx context.Context, opt online.GetByGroupOpt) (dbutil.RowMap, error) {

	if err := opt.Validate(); err != nil {
		return nil, err
	}

	key, err := serializeRedisKey(opt.Group, opt.EntityKey, opt.RevisionID)
	if err != nil {
		return nil, err
	}

	values, err := db.HGetAll(ctx, key).Result()
	if err != nil {
		return nil, errdefs.WithStack(err)
	}

	rowMap := make(dbutil.RowMap)
	for k, v := range values {
		featureID, err := dbutil.DeserializeByValueType(k, types.Int64, Backend)
		if err != nil {
			return nil, err
		}
		featureIDInt := int(featureID.(int64))
		feature, err := opt.GetFeature(&featureIDInt, nil)
		if err != nil {
			return nil, err
		}
		deserializedValue, err := dbutil.DeserializeByValueType(v, feature.ValueType, Backend)
		if err != nil {
			return nil, err
		}
		rowMap[feature.FullName()] = deserializedValue
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

func serializeRedisKey(group types.Group, entityKey string, revisionID *int) (string, error) {
	if group.Category == types.CategoryBatch {
		return serializeRedisKeyForBatchFeature(*revisionID, entityKey)
	} else {
		return serializeRedisKeyForStreamFeature(group.ID, entityKey)
	}
}
