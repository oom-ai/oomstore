package redis

import (
	"context"

	"github.com/oom-ai/oomstore/internal/database/dbutil"
	"github.com/oom-ai/oomstore/internal/database/online"
)

func (db *DB) Push(ctx context.Context, opt online.PushOpt) error {
	key, err := serializeRedisKeyForStreamFeature(opt.GroupID, opt.EntityKey)
	if err != nil {
		return err
	}

	featureValues := make(map[string]string)
	for i, value := range opt.FeatureValues {
		// omit nil feature value
		if value == nil {
			continue
		}

		featureValue, err := dbutil.SerializeByValueType(value, opt.Features[i].ValueType, Backend)
		if err != nil {
			return err
		}

		featureID, err := dbutil.SerializeByValue(opt.Features[i].ID, Backend)
		if err != nil {
			return err
		}
		featureValues[featureID] = featureValue.(string)
	}

	db.HSet(ctx, key, featureValues)
	return nil
}
